package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type AppTracer struct {
	tracer trace.Tracer
}

func NewAppTracer(serviceName string) *AppTracer {
	return &AppTracer{
		tracer: otel.Tracer(serviceName),
	}
}

func (at *AppTracer) StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return at.tracer.Start(ctx, name)
}
