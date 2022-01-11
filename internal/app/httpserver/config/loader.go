package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type loader struct {
	configPath string
}

func NewLoader(configPath string) *loader {
	return &loader{
		configPath: configPath,
	}
}

func (configLoader *loader) Load() (*Config, error) {
	cfg := NewConfig()

	err := cleanenv.ReadConfig(configLoader.configPath, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
