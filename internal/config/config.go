package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	AppPort     string
	LogLevel    string
	AutoMigrate bool
}

func Load() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "subscriptions")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("AUTO_MIGRATE", false)

	cfg := &Config{
		DBHost:     v.GetString("DB_HOST"),
		DBPort:     v.GetInt("DB_PORT"),
		DBUser:     v.GetString("DB_USER"),
		DBPassword: v.GetString("DB_PASSWORD"),
		DBName:     v.GetString("DB_NAME"),
		DBSSLMode:  v.GetString("DB_SSLMODE"),

		AppPort:     v.GetString("APP_PORT"),
		LogLevel:    v.GetString("LOG_LEVEL"),
		AutoMigrate: v.GetBool("AUTO_MIGRATE"),
	}

	if cfg.DBHost == "" || cfg.DBUser == "" {
		return nil, fmt.Errorf("invalid db config")
	}
	return cfg, nil
}
