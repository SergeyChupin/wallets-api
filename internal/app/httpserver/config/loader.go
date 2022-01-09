package config

import "github.com/SergeyChupin/wallets-api/internal/app/pkg/config"

type loader struct {
	configPath string
}

func NewLoader(configPath string) *loader {
	return &loader{
		configPath: configPath,
	}
}

func (configLoader *loader) Load() (*Config, error) {
	if err := config.Load(configLoader.configPath); err != nil {
		return nil, err
	}
	configuration := NewConfig()
	if err := config.Unmarshal(configuration); err != nil {
		return nil, err
	}
	return configuration, nil
}
