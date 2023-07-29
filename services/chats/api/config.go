package main

import (
	"github.com/caarlos0/env"
)

type config struct {
	MongoURL    string `env:"MONGO_URL" envDefault:"mongodb://mongodb:27017"`
	MongoDB     string `env:"MONGO_DB" envDefault:"users"`
	CountersURL string `env:"COUNTERS_URL" envDefault:"counters-api:50051"`
	RedisURL    string `env:"REDIS_URL"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
