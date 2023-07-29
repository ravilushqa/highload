package main

import (
	"github.com/caarlos0/env"
)

type config struct {
	MongoURL string `env:"MONGO_URL" envDefault:"mongodb://mongodb:27017"`
	MongoDB  string `env:"MONGO_DB" envDefault:"highload"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
