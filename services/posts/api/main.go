package main

import (
	"context"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/neonxp/rutina"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"go.uber.org/zap"

	kafkaproducerprovider "github.com/ravilushqa/highload/providers/kafka-producer"
	"github.com/ravilushqa/highload/providers/mongodb"
	redisprovider "github.com/ravilushqa/highload/providers/redis"
	"github.com/ravilushqa/highload/services/posts/lib/post"
)

func buildContainer() (*dig.Container, error) {
	container := dig.New()
	constructors := []interface{}{
		newConfig,
		zap.NewDevelopment,
		func(c *config) (*redis.Client, error) {
			return redisprovider.New(c.RedisURL)
		},
		func(c *config) (*mongo.Database, error) {
			return mongodb.New(context.Background(), c.MongoURL, c.MongoDB)
		},
		func(c *config) (*kafkaproducerprovider.KafkaProducer, error) {
			return kafkaproducerprovider.New(c.KafkaBrokers, c.KafkaTopic, nil)
		},
		post.NewManager,
		NewApi,
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

	err = container.Invoke(func(l *zap.Logger, r *redis.Client, db *mongo.Database) error {
		if err := r.ShutdownSave().Err(); err != nil {
			l.Error("failed to shutdown redis", zap.Error(err))
		}
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
