package main

import (
	"log"
	"os"

	"catalog-notification/internal/broker"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}

	consumer, err := broker.NewConsumer(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialize consumer: %v", err)
	}
	defer consumer.Close()

	log.Println("Notifications Service started...")
	if err := consumer.Start(); err != nil {
		log.Fatalf("Error reading from queue: %v", err)
	}
}
