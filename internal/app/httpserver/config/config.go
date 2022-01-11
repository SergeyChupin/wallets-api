package config

import (
	"github.com/SergeyChupin/wallets-api/internal/database/postgres"
	"github.com/SergeyChupin/wallets-api/internal/server"
)

type Config struct {
	Server   server.Config   `yaml:"server"`
	Postgres postgres.Config `yaml:"postgres"`
}

func NewConfig() *Config {
	return &Config{
		Server:   server.NewConfig(),
		Postgres: postgres.NewConfig(),
	}
}
