package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func Load(configPath string) error {
	if _, err := os.Stat(configPath); err != nil {
		return err
	}
	viper.SetConfigFile(configPath)
	if err := viper.MergeInConfig(); err != nil {
		return err
	}
	return nil
}

func Unmarshal(rawType interface{}) error {
	if err := viper.Unmarshal(rawType); err != nil {
		return err
	}
	return nil
}
