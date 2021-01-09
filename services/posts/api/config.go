package main

import (
	"github.com/caarlos0/env"
)

type config struct {
	DatabaseURL  string   `env:"DATABASE_URL"`
	SlavesUrls   []string `env:"SLAVES_URLS" envSeparator:","`
	KafkaBrokers []string `env:"KAFKA_BROKERS" envSeparator:","`
	KafkaTopic   string   `env:"KAFKA_TOPIC"`
	RedisURL     string   `env:"REDIS_URL"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
