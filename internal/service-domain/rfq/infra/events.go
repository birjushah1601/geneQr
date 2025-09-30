package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

// KafkaEventPublisher implements the domain.EventPublisher interface
type KafkaEventPublisher struct {
	writer *kafka.Writer
	logger *slog.Logger
}

// NewKafkaEventPublisher creates a new Kafka event publisher
func NewKafkaEventPublisher(brokers []string, logger *slog.Logger) *KafkaEventPublisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "rfq-events",
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaEventPublisher{
		writer: writer,
		logger: logger.With(slog.String("component", "kafka_event_publisher")),
	}
}

// Publish publishes a domain event to Kafka
func (p *KafkaEventPublisher) Publish(ctx context.Context, event interface{}) error {
	// Serialize event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal event",
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to Kafka
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Value: eventJSON,
	})

	if err != nil {
		p.logger.Error("Failed to publish event to Kafka",
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Debug("Event published successfully")
	return nil
}

// Close closes the Kafka writer
func (p *KafkaEventPublisher) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
