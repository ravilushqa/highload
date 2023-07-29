package main

import (
	"context"

	"github.com/axengine/go-saga"
	_ "github.com/go-sql-driver/mysql"
	"github.com/neonxp/rutina"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/ravilushqa/highload/providers/mongodb"
	sagaProvider "github.com/ravilushqa/highload/providers/saga"
	"github.com/ravilushqa/highload/services/chats/lib/chat"
	"github.com/ravilushqa/highload/services/chats/lib/chat_user"
	"github.com/ravilushqa/highload/services/chats/lib/message"
	countersGrpc "github.com/ravilushqa/highload/services/counters/api/grpc"
)

func buildContainer() (*dig.Container, error) {
	container := dig.New()
	constructors := []interface{}{
		newConfig,
		zap.NewDevelopment,
		func(c *config) (*mongo.Database, error) {
			return mongodb.New(context.Background(), c.MongoURL, c.MongoDB)
		},
		message.NewManager,
		NewApi,
		chat.NewManager,
		chatuser.NewManager,
		func(c *config) (countersGrpc.CountersClient, error) {
			conn, err := grpc.Dial(c.CountersURL, grpc.WithInsecure())
			if err != nil {
				return nil, err
			}
			return countersGrpc.NewCountersClient(conn), nil
		},
		func(c *config) (*saga.ExecutionCoordinator, error) {
			return sagaProvider.New(c.RedisURL)
		},
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

	err = container.Invoke(func(l *zap.Logger, db *mongo.Database) error {
		if err := db.Client().Disconnect(context.Background()); err != nil {
			l.Error("failed disconnect from mongo", zap.Error(err))
		}
		l.Info("gracefully shutdown...")
		return l.Sync()
	})
	if err != nil {
		tl.Error("shutdown failed", zap.Error(err))
	}

}
