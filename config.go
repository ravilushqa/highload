package main

import (
	"time"

	"github.com/caarlos0/env"
)

type config struct {
	LogLevel         string        `env:"LOG_LEVEL" envDefault:"debug"`
	Addr             string        `env:"ADDR" envDefault:":8080"`
	DatabaseURL      string        `env:"DATABASE_URL"`
	SlavesUrls       []string      `env:"SLAVES_URLS" envSeparator:","`
	MessagesShards   []string      `env:"MESSAGES_SHARDS" envSeparator:","`
	JwtSecret        string        `env:"JWT_SECRET" envDefault:"secret"`
	APITimeout       time.Duration `env:"API_TIMEOUT" envDefault:"60s"`
	KafkaBrokers     []string      `env:"KAFKA_BROKERS" envSeparator:","`
	KafkaTopic       string        `env:"KAFKA_TOPIC"`
	KafkaGroupID     string        `env:"KAFKA_GROUP_ID" envDefault:"app"`
	RedisURL         string        `env:"REDIS_URL"`
	TarantoolURL     string        `env:"TARANTOOL_URL"`
	TarantoolUser    string        `env:"TARANTOOL_USER" envDefault:"guest"`
	TarantoolPass    string        `env:"TARANTOOL_PASS"`
	CentrifugoURL    string        `env:"CENTRIFUGO_URL" envDefault:"http://centrifugo:8000"`
	CentrifugoApiKey string        `env:"CENTRIFUGO_API_KEY" envDefault:"my_api_key"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
