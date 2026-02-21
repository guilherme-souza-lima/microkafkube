package queue

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaAdapterInterface interface {
	Publish(ctx context.Context, key string, payload []byte) error
}

type KafkaAdapter struct {
	writer *kafka.Writer
}

func NewKafkaAdapter(writer *kafka.Writer) *KafkaAdapter {
	return &KafkaAdapter{writer: writer}
}

func (a *KafkaAdapter) Publish(ctx context.Context, key string, payload []byte) error {
	return a.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: payload,
	})
}
