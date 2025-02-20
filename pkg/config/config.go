package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port         int    `mapstructure:"PORT"`
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	ExamPort     int    `mapstructure:"EXAM_PORT"`
	QuestionPort int    `mapstructure:"QUESTION_PORT"`
	SessionPort  int    `mapstructure:"SESSION_PORT"`
	ScoringPort  int    `mapstructure:"SCORING_PORT"`
}

func Load() (*Config, error) {
	// Default values
	viper.SetDefault("PORT", 50051)
	viper.SetDefault("EXAM_PORT", 50052)
	viper.SetDefault("QUESTION_PORT", 50053)
	viper.SetDefault("SESSION_PORT", 50054)
	viper.SetDefault("SCORING_PORT", 50055)
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/cbt_exam?sslmode=disable")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
