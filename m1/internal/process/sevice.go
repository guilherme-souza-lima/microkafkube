package process

import (
	"context"
	"fmt"
	"microum/queue"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

type ServiceInterface interface {
	Register(ctx context.Context, data []byte) (ResponseRegisterDTO, error)
}

type Service struct {
	repository             RepositoryInterface
	regCounter, errCounter prometheus.Counter
	messaging              queue.KafkaAdapterInterface
}

func NewService(repository RepositoryInterface, counter, errCounter prometheus.Counter, messaging queue.KafkaAdapterInterface) *Service {
	return &Service{repository, counter, errCounter, messaging}
}

func (s *Service) Register(ctx context.Context, data []byte) (ResponseRegisterDTO, error) {
	traceID := uuid.New()
	lenBytes := len(data)
	lenCharacters := utf8.RuneCount(data)

	dto := RegisterDTO{
		TraceID:         traceID,
		Payload:         data,
		ByteSize:        lenBytes,
		TotalCharacters: lenCharacters,
	}

	err := s.repository.Register(ctx, dto)
	if err != nil {
		return ResponseRegisterDTO{}, fmt.Errorf("service failed to register: %w", err)
	}

	s.regCounter.Inc()

	go func(tID string, payload []byte) {
		kafkaCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.messaging.Publish(kafkaCtx, tID, payload)
		if err == nil {
			_ = s.repository.UpdatePublishedStatus(context.Background(), tID)
		} else {
			s.errCounter.Inc()
			fmt.Printf("⚠️ Background Queue delivery failed for %s: %v\n", tID, err)
		}
	}(traceID.String(), data)

	return ResponseRegisterDTO{
		TraceID: traceID.String(),
	}, nil
}
