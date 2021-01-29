package saga

import (
	"github.com/axengine/go-saga"
	"github.com/axengine/go-saga/storage/redis"
)

const prefix = "saga"

func New(redisURL string) (*saga.ExecutionCoordinator, error) {
	store, err := redis.NewRedisStore(
		redisURL,
		"",
		1,
		2,
		5,
		prefix,
	)
	if err != nil {
		return nil, err
	}
	sec := saga.NewSEC(store, prefix)

	return &sec, nil
}
