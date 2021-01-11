package main

import (
	"github.com/caarlos0/env"
)

type config struct {
	KafkaBrokers     []string `env:"KAFKA_BROKERS" envSeparator:","`
	KafkaTopic       string   `env:"KAFKA_TOPIC"`
	KafkaGroupID     string   `env:"KAFKA_GROUP_ID" envDefault:"app"`
	RedisURL         string   `env:"REDIS_URL"`
	CentrifugoURL    string   `env:"CENTRIFUGO_URL" envDefault:"http://centrifugo:8000"`
	CentrifugoApiKey string   `env:"CENTRIFUGO_API_KEY" envDefault:"my_api_key"`
	UsersURL         string   `env:"USERS_URL" envDefault:"users-api:50051"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
