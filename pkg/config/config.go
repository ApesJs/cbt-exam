package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port        int    `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

func Load() (*Config, error) {
	viper.SetDefault("PORT", 50051)
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/cbt_exam?sslmode=disable")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
