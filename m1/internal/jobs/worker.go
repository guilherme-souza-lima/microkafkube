package jobs

import (
	"context"
	"fmt"
	"microum/internal/process"
	"microum/internal/queue"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type OutboxJob struct {
	repository process.OutboxRepository
	kafka      queue.KafkaAdapterInterface
	errCounter prometheus.Counter
}

func NewOutboxJob(repo process.OutboxRepository, kafka queue.KafkaAdapterInterface, errCounter prometheus.Counter) *OutboxJob {
	return &OutboxJob{
		repository: repo,
		kafka:      kafka,
		errCounter: errCounter,
	}
}

func (w *OutboxJob) Start(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processOutbox(ctx)
		}
	}
}

func (w *OutboxJob) processOutbox(ctx context.Context) {
	pendings, _ := w.repository.GetPendingRegistrations(ctx)

	for _, reg := range pendings {
		err := w.kafka.Publish(ctx, reg.TraceID, reg.Payload)
		if err == nil {
			_ = w.repository.UpdatePublishedStatus(ctx, reg.TraceID)
		} else {
			w.errCounter.Inc()
			fmt.Printf("⚠️OutboxJob: Queue delivery failed for %s: %v\n", reg.TraceID, err)
		}
	}
}
