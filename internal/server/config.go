package server

import "time"

type Config struct {
	Port                string        `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout         time.Duration `yaml:"read-timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout        time.Duration `yaml:"write-timeout" env:"SERVER_WRITE_TIMEOUT"`
	ShutdownGracePeriod time.Duration `yaml:"shutdown-grace-period" env:"SERVER_SHUTDOWN_GRACE_PERIOD"`
}

func NewConfig() Config {
	return Config{
		Port:                "8080",
		ReadTimeout:         time.Second * 5,
		WriteTimeout:        time.Second * 15,
		ShutdownGracePeriod: time.Second * 30,
	}
}
