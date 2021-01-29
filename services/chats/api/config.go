package main

import (
	"github.com/caarlos0/env"
)

type config struct {
	DatabaseURL    string   `env:"DATABASE_URL"`
	SlavesUrls     []string `env:"SLAVES_URLS" envSeparator:","`
	MessagesShards []string `env:"MESSAGES_SHARDS" envSeparator:","`
	CountersURL    string   `env:"COUNTERS_URL" envDefault:"counters-api:50051"`
	RedisURL       string   `env:"REDIS_URL"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
