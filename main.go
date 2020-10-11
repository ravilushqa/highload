package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/neonxp/rutina"
	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/ravilushqa/highload/controllers/auth"
	"github.com/ravilushqa/highload/controllers/users"
	"github.com/ravilushqa/highload/lib"
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
		func(c *config) (*sqlx.DB, error) {
			return sqlx.Connect("mysql", fmt.Sprint(c.MysqlConnection, "?parseTime=true"))
		},
		NewAPI,
		user.New,
		auth.NewController,
		users.NewController,
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
			tl.Error("runtime error", zap.NamedError("b2bError", err))
		}
	}()

	err = container.Invoke(func(api *API) {
		r.Go(api.Run)
		r.ListenOsSignals()
	})
	if err != nil {
		tl.Fatal("invoke failed", zap.NamedError("b2bError", err))
	}

	if err := r.Wait(); err != nil {
		tl.Error("run failed ", zap.NamedError("b2bError", err))
	}

	err = container.Invoke(func(l *zap.Logger) error {
		l.Info("gracefully shutdown...")
		return l.Sync()
	})
	if err != nil {
		tl.Error("shutdown failed", zap.NamedError("b2bError", err))
	}

}
