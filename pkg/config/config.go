package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port         int    `mapstructure:"PORT"`
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	ExamPort     int    `mapstructure:"EXAM_PORT"`
	QuestionPort int    `mapstructure:"QUESTION_PORT"`
	SessionPort  int    `mapstructure:"SESSION_PORT"`
	ScoringPort  int    `mapstructure:"SCORING_PORT"`
	DatabaseHost string `mapstructure:"DB_HOST"`
	DatabasePort int    `mapstructure:"DB_PORT"`
	DatabaseName string `mapstructure:"DB_NAME"`
	DatabaseUser string `mapstructure:"DB_USER"`
	DatabasePass string `mapstructure:"DB_PASS"`
	DatabaseSSL  string `mapstructure:"DB_SSL"`
}

func Load() (*Config, error) {
	// Set default values
	//setDefaults()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
		// Config file not found, will use defaults and env vars
		fmt.Println("No .env file found, using environment variables or defaults")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %v", err)
	}

	// Build database URL if not set
	if config.DatabaseURL == "" {
		config.DatabaseURL = buildDatabaseURL(config)
	}

	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

//func setDefaults() {
//	// Service ports
//	//viper.SetDefault("PORT", 8080)           // Default main port
//	viper.SetDefault("EXAM_PORT", 50051)     // Exam service port
//	viper.SetDefault("QUESTION_PORT", 50052) // Question service port
//	viper.SetDefault("SESSION_PORT", 50053)  // Session service port
//	viper.SetDefault("SCORING_PORT", 50054)  // Scoring service port
//
//	// Database defaults
//	viper.SetDefault("DB_HOST", "localhost")
//	viper.SetDefault("DB_PORT", 5432)
//	viper.SetDefault("DB_NAME", "cbt_exam")
//	viper.SetDefault("DB_USER", "root")
//	viper.SetDefault("DB_PASS", "secret")
//	viper.SetDefault("DB_SSL", "disable")
//
//	// Default database URL (will be overridden if individual components are set)
//	viper.SetDefault("DATABASE_URL", "postgres://root:secret@localhost:5432/cbt_exam?sslmode=disable")
//
//}

func buildDatabaseURL(config Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.DatabaseUser,
		config.DatabasePass,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
		config.DatabaseSSL,
	)
}

func validateConfig(config *Config) error {
	// Validate required ports
	if config.Port <= 0 {
		return fmt.Errorf("invalid port number: %d", config.Port)
	}

	// Validate service ports
	servicePorts := []struct {
		name string
		port int
	}{
		{"EXAM", config.ExamPort},
		{"QUESTION", config.QuestionPort},
		{"SESSION", config.SessionPort},
		{"SCORING", config.ScoringPort},
	}

	for _, sp := range servicePorts {
		if sp.port <= 0 {
			return fmt.Errorf("invalid %s service port number: %d", sp.name, sp.port)
		}
	}

	// Validate database configuration
	if config.DatabaseURL == "" {
		return fmt.Errorf("database URL is required")
	}

	return nil
}

// Helper function untuk mendapatkan environment variable dengan default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
