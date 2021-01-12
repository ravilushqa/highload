package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/linxGnu/mssqlx"
	"github.com/neonxp/rutina"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/providers/db"
	tarantoolprovider "github.com/ravilushqa/highload/providers/tarantool"
	"github.com/ravilushqa/highload/services/users/lib/friend"
	"github.com/ravilushqa/highload/services/users/lib/user"
)

func buildContainer() (*dig.Container, error) {
	container := dig.New()
	constructors := []interface{}{
		newConfig,
		zap.NewDevelopment,
		func(c *config) (*mssqlx.DBs, error) {
			return db.New(c.DatabaseURL, c.SlavesUrls)
		},
		func(c *config) (*tarantool.Connection, error) {
			return tarantoolprovider.New(c.TarantoolURL, c.TarantoolUser, c.TarantoolPass)
		},
		NewApi,
		user.New,
		friend.New,
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
