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

	"github.com/ravilushqa/highload/services/chats/lib/chat"
	"github.com/ravilushqa/highload/services/chats/lib/chat_user"
	"github.com/ravilushqa/highload/services/chats/lib/message"

	"github.com/ravilushqa/highload/providers/db"
)

func buildContainer() (*dig.Container, error) {
	container := dig.New()
	constructors := []interface{}{
		newConfig,
		zap.NewDevelopment,
		func(c *config) (*mssqlx.DBs, error) {
			return db.New(c.DatabaseURL, c.SlavesUrls)
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
		NewApi,
		chat.NewManager,
		chatuser.NewManager,
	}

	for _, c := range constructors {
		if err := container.Provide(c); err != nil {
			return nil, err
		}
	}

	return container, container.Invoke(func(a *Api) {})
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

	err = container.Invoke(func(api *Api) {
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
