package main

import (
	"time"

	"github.com/caarlos0/env"
)

type config struct {
	Addr       string        `env:"ADDR" envDefault:":8080"`
	JwtSecret  string        `env:"JWT_SECRET" envDefault:"secret"`
	APITimeout time.Duration `env:"API_TIMEOUT" envDefault:"60s"`
	ChatsURL   string        `env:"CHATS_URL" envDefault:"chats-api:50051"`
	PostsURL   string        `env:"POSTS_URL" envDefault:"posts-api:50051"`
	UsersURL   string        `env:"USERS_URL" envDefault:"users-api:50051"`
}

func newConfig() (*config, error) {
	cfg := new(config)
	return cfg, env.Parse(cfg)
}
