package infra

import (
	"context"
	"log"
	"microdois/internal/messaging"
	"microdois/internal/observability"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Container struct {
	Config         *Config
	MongoDatabase  *mongo.Database
	KafkaConsumer  *messaging.KafkaConsumer
	TracerProvider *trace.TracerProvider
	Tracer         *observability.AppTracer
	KafkaCounter   prometheus.Counter
}

func NewContainer(config *Config, ctx context.Context) *Container {
	container := &Container{}
	container.Config = config

	container.buildTracer(ctx)
	container.buildPrometheus()
	container.buildMongoDB(ctx)
	container.buildKafkaConsumer()

	return container
}

func (c *Container) buildTracer(ctx context.Context) {}
func (c *Container) buildPrometheus()                {}
func (c *Container) buildMongoDB(ctx context.Context) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctxTimeout, options.Client().ApplyURI(c.Config.MongoURI))
	if err != nil {
		log.Fatalf("❌ M2: Erro ao conectar no MongoDB: %v", err)
	}

	if err := client.Ping(ctxTimeout, nil); err != nil {
		log.Fatalf("❌ M2: MongoDB não responde: %v", err)
	}

	c.MongoDatabase = client.Database(c.Config.MongoDatabase)
	log.Println("✅ M2: MongoDB conectado com sucesso")
}
func (c *Container) buildKafkaConsumer() {}
