package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"catalog-product/internal/model"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "catalog_events"
)

type RabbitMQBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQBroker(url string) (*RabbitMQBroker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	slog.Info("Connected to RabbitMQ successfully")

	return &RabbitMQBroker{
		conn:    conn,
		channel: ch,
	}, nil
}

func (b *RabbitMQBroker) Close() {
	if b.channel != nil {
		b.channel.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}
}

func (b *RabbitMQBroker) publishEvent(ctx context.Context, action string, data interface{}) error {
	message := map[string]interface{}{
		"action": action,
		"data":   data,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	err = b.channel.PublishWithContext(ctx,
		exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to RabbitMQ: %w", err)
	}

	slog.Info("Published event to RabbitMQ", "action", action)
	return nil
}

func (b *RabbitMQBroker) PublishProductCreated(ctx context.Context, product *model.Product) error {
	return b.publishEvent(ctx, "product.created", product)
}

func (b *RabbitMQBroker) PublishProductDeleted(ctx context.Context, id uuid.UUID) error {
	return b.publishEvent(ctx, "product.deleted", map[string]string{
		"product_id": id.String(),
	})
}
