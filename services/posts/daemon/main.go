package main

import (
	cluster "github.com/bsm/sarama-cluster"
	"github.com/centrifugal/gocent"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/linxGnu/mssqlx"
	"github.com/neonxp/rutina"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	centrifugoclient "github.com/ravilushqa/highload/providers/centrifugo-client"
	kafkaconsumerprovider "github.com/ravilushqa/highload/providers/kafka-consumer"
	redisprovider "github.com/ravilushqa/highload/providers/redis"
	usersGrpc "github.com/ravilushqa/highload/services/users/api/grpc"
)

func buildContainer() (*dig.Container, error) {
	container := dig.New()
	constructors := []interface{}{
		newConfig,
		zap.NewDevelopment,
		func(c *config) (*redis.Client, error) {
			return redisprovider.New(c.RedisURL)
		},
		func(c *config) (*cluster.Consumer, error) {
			return kafkaconsumerprovider.New(c.KafkaBrokers, c.KafkaGroupID, []string{c.KafkaTopic}, nil)
		},
		func(c *config) *gocent.Client {
			return centrifugoclient.New(c.CentrifugoURL, c.CentrifugoApiKey)
		},
		newDaemon,
		func(c *config) (usersGrpc.UsersClient, error) {
			conn, err := grpc.Dial(c.UsersURL, grpc.WithInsecure())
			if err != nil {
				return nil, err
			}
			return usersGrpc.NewUsersClient(conn), nil
		},
	}

	for _, c := range constructors {
		if err := container.Provide(c); err != nil {
			return nil, err
		}
	}

	return container, container.Invoke(func(d *daemon) {})
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

	err = container.Invoke(func(d *daemon) {
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
