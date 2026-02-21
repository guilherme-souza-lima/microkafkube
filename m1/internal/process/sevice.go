package process

import (
	"context"
	"fmt"
	"microum/internal/observability"
	"microum/queue"
	"time"
	"unicode/utf8"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
)

type ServiceInterface interface {
	Register(ctx context.Context, data []byte) (ResponseRegisterDTO, error)
}

type Service struct {
	repository             RepositoryInterface
	regCounter, errCounter prometheus.Counter
	messaging              queue.KafkaAdapterInterface
	tracer                 observability.TracerInterface
}

func NewService(
	repository RepositoryInterface,
	counter prometheus.Counter,
	errCounter prometheus.Counter,
	messaging queue.KafkaAdapterInterface,
	tracer observability.TracerInterface,
) *Service {
	return &Service{
		repository: repository,
		regCounter: counter,
		errCounter: errCounter,
		messaging:  messaging,
		tracer:     tracer,
	}
}

func (s *Service) Register(ctx context.Context, data []byte) (ResponseRegisterDTO, error) {
	ctx, span := s.tracer.StartSpan(ctx, "Service.Register")
	defer span.End()

	tID := span.SpanContext().TraceID().String()

	ctxDB, dbSpan := s.tracer.StartSpan(ctx, "Repository.Save")
	dto := RegisterDTO{
		TraceID:         tID,
		Payload:         data,
		ByteSize:        len(data),
		TotalCharacters: utf8.RuneCount(data),
	}
	err := s.repository.Register(ctxDB, dto)
	dbSpan.End()

	if err != nil {
		span.RecordError(err)
		return ResponseRegisterDTO{}, fmt.Errorf("service failed to register: %w", err)
	}

	s.regCounter.Inc()

	go func(currentCtx context.Context, id string, payload []byte) {
		kafkaCtx := trace.ContextWithSpanContext(context.Background(), trace.SpanContextFromContext(currentCtx))
		kafkaCtx, cancel := context.WithTimeout(kafkaCtx, 5*time.Second)
		defer cancel()

		kafkaCtx, kSpan := s.tracer.StartSpan(kafkaCtx, "Kafka.BackgroundPublish")
		defer kSpan.End()

		err := s.messaging.Publish(kafkaCtx, id, payload)
		if err == nil {
			_ = s.repository.UpdatePublishedStatus(context.Background(), id)
		} else {
			kSpan.RecordError(err)
			s.errCounter.Inc()
			fmt.Printf("⚠️ Background Queue delivery failed for %s: %v\n", id, err)
		}
	}(ctx, tID, data)

	return ResponseRegisterDTO{
		TraceID: tID,
	}, nil
}
