package server

import "time"

type Config struct {
	Port                string        `mapstructure:"port"`
	ReadTimeout         time.Duration `mapstructure:"read-timeout"`
	WriteTimeout        time.Duration `mapstructure:"write-timeout"`
	ShutdownGracePeriod time.Duration `mapstructure:"shutdown-grace-period"`
}

func NewConfig() Config {
	return Config{
		Port:                "8080",
		ReadTimeout:         time.Second * 5,
		WriteTimeout:        time.Second * 15,
		ShutdownGracePeriod: time.Second * 30,
	}
}
