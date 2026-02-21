package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// TracerInterface define o que a Service e o Handler podem fazer
type TracerInterface interface {
	StartSpan(ctx context.Context, name string) (context.Context, trace.Span)
}

type AppTracer struct {
	tracer trace.Tracer
}

// NewAppTracer cria a instância que será injetada nos componentes
func NewAppTracer(serviceName string) *AppTracer {
	return &AppTracer{
		tracer: otel.Tracer(serviceName),
	}
}

// StartSpan inicia um novo segmento de tempo no rastro
func (at *AppTracer) StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return at.tracer.Start(ctx, name)
}

// InitGlobalTracer configura o exportador que enviará os dados para o Jaeger/OTel Collector
func InitGlobalTracer(ctx context.Context, serviceName, endpoint string) (*sdktrace.TracerProvider, error) {
	// 1. Configura o exportador via gRPC (Porta 4317 padrão)
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// 2. Define os atributos do serviço (aparecerão no Jaeger)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 3. Cria o provedor do rastro
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// 4. DEFINE OS GLOBAIS (O segredo da propagação)
	otel.SetTracerProvider(tp)

	// Isso aqui permite que o rastro viaje nos Headers (HTTP ou Kafka)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}
