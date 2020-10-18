package main

import (
	"time"

	"github.com/caarlos0/env"
)

type config struct {
	LogLevel    string        `env:"LOG_LEVEL" envDefault:"debug"`
	Addr        string        `env:"ADDR" envDefault:":8080"`
	DatabaseURL string        `env:"DATABASE_URL" envDefault:"user:password@(localhost:3306)/app"`
	JwtSecret   string        `env:"JWT_SECRET" envDefault:"secret"`
	APITimeout  time.Duration `env:"API_TIMEOUT" envDefault:"60s"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
