package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
	"github.com/neonxp/rutina"
	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/controllers/auth"
	"github.com/ravilushqa/highload/controllers/chats"
	"github.com/ravilushqa/highload/controllers/users"
	"github.com/ravilushqa/highload/lib"
	"github.com/ravilushqa/highload/lib/chat"
	chatuser "github.com/ravilushqa/highload/lib/chat_user"
	"github.com/ravilushqa/highload/lib/friend"
	"github.com/ravilushqa/highload/lib/message"
	"github.com/ravilushqa/highload/lib/user"
)

func buildContainer() (*dig.Container, error) {
	container := dig.New()
	constructors := []interface{}{
		newConfig,
		zap.NewProduction,
		func(c *config) (*lib.Auth, error) {
			return lib.NewAuth(c.JwtSecret), nil
		},
		func(c *config) (*mssqlx.DBs, error) {
			dsns := make([]string, 0, len(c.SlavesUrls)+1)
			dsns = append(dsns, c.DatabaseURL)
			dsns = append(dsns, c.SlavesUrls...)
			for i := range dsns {
				dsns[i] = dsns[i] + "?parseTime=true"
			}
			db, errs := mssqlx.ConnectMasterSlaves("mysql", dsns[:1], dsns[1:])

			for _, err := range errs {
				if err != nil {
					return nil, fmt.Errorf("failed init db connection: %v", errs)
				}
			}

			//db.SetMaxOpenConns(25)
			//db.SetMaxIdleConns(25)
			db.SetConnMaxLifetime(5 * time.Minute)
			errs = db.Ping()
			for _, err := range errs {
				if err != nil {
					return nil, fmt.Errorf("database is unreachable: %v", errs)
				}
			}

			return db, nil
		},
		func(c *config) (*message.Manager, error) {
			dbs := make([]*sqlx.DB, 0)
			r := rutina.New()
			for _, shardURL := range c.MessagesShards {
				r.Go(func(ctx context.Context) error {
					db, err := sqlx.Connect("mysql", fmt.Sprint(shardURL, "?parseTime=true"))
					if err != nil {
						return err
					}

					db.SetConnMaxLifetime(5 * time.Minute)
					dbs = append(dbs, db)
					return nil
				})
			}

			err := r.Wait()
			if err != nil {
				return nil, err
			}
			return message.NewManager(dbs), nil
		},
		NewAPI,
		user.New,
		friend.New,
		auth.NewController,
		users.NewController,
		chat.NewManager,
		chatuser.NewManager,
		chats.NewController,
	}

	for _, c := range constructors {
		if err := container.Provide(c); err != nil {
			return nil, err
		}
	}

	return container, container.Invoke(func(a *API) {})
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

	err = container.Invoke(func(api *API) {
		r.Go(api.Run)
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
