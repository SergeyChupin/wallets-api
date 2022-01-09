package config

import (
	"github.com/SergeyChupin/wallets-api/internal/app/pkg/database/postgres"
	"github.com/SergeyChupin/wallets-api/internal/app/pkg/server"
)

type Config struct {
	Server   server.Config   `mapstructure:"server"`
	Postgres postgres.Config `mapstructure:"postgres"`
}

func NewConfig() *Config {
	return &Config{
		Server:   server.NewConfig(),
		Postgres: postgres.NewConfig(),
	}
}
