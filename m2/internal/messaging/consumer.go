package messaging

import (
	"microdois/internal/observability"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	tracer *observability.AppTracer
}

func NewKafkaConsumer(addr, topic, groupID string, tracer *observability.AppTracer) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{addr},
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		tracer: tracer,
	}
}
