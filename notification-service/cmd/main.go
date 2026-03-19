package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"catalog-notification/internal/broker"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		slog.Error("RABBITMQ_URL is not set")
		os.Exit(1)
	}

	consumer, err := broker.NewConsumer(rabbitURL)
	if err != nil {
		slog.Error("Failed to initialize consumer", "error", err)
		os.Exit(1)
	}
	defer consumer.Close()

	go func() {
		slog.Info("Notifications Service started...")
		if err := consumer.Start(); err != nil {
			slog.Error("Error reading from queue", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down Notifications Service...")
}
