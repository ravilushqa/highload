package main

import (
	"context"
	"fmt"
	"time"

	cluster "github.com/bsm/sarama-cluster"
	"github.com/centrifugal/gocent"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
	"github.com/neonxp/rutina"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/controllers/auth"
	"github.com/ravilushqa/highload/controllers/chats"
	"github.com/ravilushqa/highload/controllers/feed"
	"github.com/ravilushqa/highload/controllers/posts"
	"github.com/ravilushqa/highload/controllers/users"
	"github.com/ravilushqa/highload/lib"
	"github.com/ravilushqa/highload/lib/chat"
	chatuser "github.com/ravilushqa/highload/lib/chat_user"
	"github.com/ravilushqa/highload/lib/friend"
	"github.com/ravilushqa/highload/lib/message"
	"github.com/ravilushqa/highload/lib/post"
	"github.com/ravilushqa/highload/lib/user"
	centrifugoclient "github.com/ravilushqa/highload/providers/centrifugo-client"
	"github.com/ravilushqa/highload/providers/db"
	kafkaconsumerprovider "github.com/ravilushqa/highload/providers/kafka-consumer"
	kafkaproducerprovider "github.com/ravilushqa/highload/providers/kafka-producer"
	redisprovider "github.com/ravilushqa/highload/providers/redis"
	tarantoolprovider "github.com/ravilushqa/highload/providers/tarantool"
)

func buildContainer() (*dig.Container, error) {
	container := dig.New()
	constructors := []interface{}{
		newConfig,
		zap.NewDevelopment,
		func(c *config) (*lib.Auth, error) {
			return lib.NewAuth(c.JwtSecret), nil
		},
		func(c *config) (*mssqlx.DBs, error) {
			return db.New(c.DatabaseURL, c.SlavesUrls)
		},
		func(c *config) (*redis.Client, error) {
			return redisprovider.New(c.RedisURL)
		},
		func(c *config) (*kafkaproducerprovider.KafkaProducer, error) {
			return kafkaproducerprovider.New(c.KafkaBrokers, c.KafkaTopic, nil)
		},
		func(c *config) (*cluster.Consumer, error) {
			return kafkaconsumerprovider.New(c.KafkaBrokers, c.KafkaGroupID, []string{c.KafkaTopic}, nil)
		},
		func(c *config) (*tarantool.Connection, error) {
			return tarantoolprovider.New(c.TarantoolURL, c.TarantoolUser, c.TarantoolPass)
		},
		func(c *config) (*message.Manager, error) {
			dbs := make([]*sqlx.DB, 0)
			r := rutina.New()
			for _, shardURL := range c.MessagesShards {
				r.Go(func(ctx context.Context) error {
					database, err := sqlx.Connect("mysql", fmt.Sprint(shardURL, "?parseTime=true"))
					if err != nil {
						return err
					}

					database.SetConnMaxLifetime(5 * time.Minute)
					database.SetConnMaxIdleTime(5 * time.Minute)
					database.SetMaxOpenConns(25)
					database.SetMaxIdleConns(25)
					dbs = append(dbs, database)
					return nil
				})
			}

			err := r.Wait()
			if err != nil {
				return nil, err
			}
			return message.NewManager(dbs), nil
		},
		func(c *config) *gocent.Client {
			return centrifugoclient.New(c.CentrifugoURL, c.CentrifugoApiKey)
		},
		NewAPI,
		newDaemon,
		user.New,
		friend.New,
		auth.NewController,
		users.NewController,
		chat.NewManager,
		chatuser.NewManager,
		post.NewManager,
		chats.NewController,
		posts.NewController,
		feed.NewController,
	}

	for _, c := range constructors {
		if err := container.Provide(c); err != nil {
			return nil, err
		}
	}

	return container, container.Invoke(func(a *API, d *daemon) {})
}

func main() {
	tl, _ := zap.NewDevelopment()
	container, err := buildContainer()
	if err != nil {
		tl.Fatal("cannot build depends", zap.Error(err))
	}

	r := rutina.New(rutina.WithErrChan())
	go func() {
		for err := range r.Errors() {
			tl.Error("runtime error", zap.Error(err))
		}
	}()

	err = container.Invoke(func(api *API, d *daemon) {
		r.Go(api.run)
		r.Go(d.run)
		r.ListenOsSignals()
	})
	if err != nil {
		tl.Fatal("invoke failed", zap.Error(err))
	}

	if err := r.Wait(); err != nil {
		tl.Error("run failed", zap.Error(err))
	}

	err = container.Invoke(func(l *zap.Logger, db *mssqlx.DBs) error {
		if errs := db.Destroy(); len(errs) > 0 {
			l.Error("failed to close db", zap.Errors("errors", errs))
		}
		l.Info("gracefully shutdown...")
		return l.Sync()
	})
	if err != nil {
		tl.Error("shutdown failed", zap.Error(err))
	}

}
