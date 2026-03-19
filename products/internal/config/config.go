package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	Port          string
	MigrationPath string
	RabbitMQURL   string
}

func LoadConfig() *Config {
	_ = godotenv.Load(".env", "../../.env", "../.env")

	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		Port:          os.Getenv("PORT"),
		MigrationPath: os.Getenv("MIGRATION_PATH"),
		RabbitMQURL:   os.Getenv("RABBITMQ_URL"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.DatabaseURL == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	if cfg.MigrationPath == "" {
		slog.Error("MIGRATION_PATH is not set")
		os.Exit(1)
	}

	if cfg.RabbitMQURL == "" {
		slog.Error("RABBITMQ_URL is not set")
		os.Exit(1)
	}

	return cfg
}
