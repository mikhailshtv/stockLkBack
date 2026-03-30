package producer

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

//go:generate mockgen -source=producer.go -destination=mocks/mock.go

type EventPublisher interface {
	Publish(ctx context.Context, key string, value any) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func New(brokers []string, topic string, logger *zap.Logger) *KafkaProducer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}
	return &KafkaProducer{writer: w, logger: logger}
}

func (p *KafkaProducer) Publish(ctx context.Context, key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: b,
	})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
