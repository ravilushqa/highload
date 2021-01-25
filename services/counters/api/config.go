package main

import (
	"github.com/caarlos0/env"
)

type config struct {
	RedisURL string `env:"REDIS_URL"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
