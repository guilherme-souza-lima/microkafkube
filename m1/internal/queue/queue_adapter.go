package queue

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	msg := kafka.Message{
		Key:   []byte(key),
		Value: payload,
	}

	carrier := propagation.HeaderCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	for k, v := range carrier {
		if len(v) > 0 {
			msg.Headers = append(msg.Headers, kafka.Header{
				Key:   k,
				Value: []byte(v[0]),
			})
		}
	}

	return a.writer.WriteMessages(ctx, msg)
}
