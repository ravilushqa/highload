package main

import (
	"github.com/caarlos0/env"
)

type config struct {
	DatabaseURL   string   `env:"DATABASE_URL"`
	SlavesUrls    []string `env:"SLAVES_URLS" envSeparator:","`
	TarantoolURL  string   `env:"TARANTOOL_URL"`
	TarantoolUser string   `env:"TARANTOOL_USER" envDefault:"guest"`
	TarantoolPass string   `env:"TARANTOOL_PASS"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
