package config

import (
	"log"
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
		log.Fatal("DATABASE_URL is not set")
	}

	if cfg.MigrationPath == "" {
		log.Fatal("MIGRATION_PATH is not set")
	}

	if cfg.RabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}

	return cfg
}
