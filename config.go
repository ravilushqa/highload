package main

import (
	"time"

	"github.com/caarlos0/env"
)

type config struct {
	LogLevel       string        `env:"LOG_LEVEL" envDefault:"debug"`
	Addr           string        `env:"ADDR" envDefault:":8080"`
	DatabaseURL    string        `env:"DATABASE_URL"`
	SlavesUrls     []string      `env:"SLAVES_URLS" envSeparator:","`
	MessagesShards []string      `env:"MESSAGES_SHARDS" envSeparator:","`
	JwtSecret      string        `env:"JWT_SECRET" envDefault:"secret"`
	APITimeout     time.Duration `env:"API_TIMEOUT" envDefault:"60s"`
	KafkaBrokers   []string      `env:"KAFKA_BROKERS" envSeparator:","`
	KafkaTopic     string        `env:"KAFKA_TOPIC"`
	KafkaGroupID   string        `env:"KAFKA_GROUP_ID" envDefault:"app"`
	RedisURL       string        `env:"REDIS_URL"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
