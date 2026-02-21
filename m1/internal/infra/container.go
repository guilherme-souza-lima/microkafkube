package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"microum"
	"microum/internal/jobs"
	"microum/internal/process"
	"microum/queue"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/segmentio/kafka-go"
)

type Container struct {
	Config     *Config
	DB         *sql.DB
	Kafka      *queue.KafkaAdapter
	RegCounter prometheus.Counter
	ErrCounter prometheus.Counter
	Repository *process.Repository
	Service    *process.Service
	Handler    *process.Handler
}

func NewContainer(config *Config, ctx context.Context) *Container {
	container := &Container{}
	container.Config = config

	container.buildDB()
	container.buildMigration()
	container.buildKafka()
	container.buildPrometheus()
	container.buildRepository()
	container.buildService()
	container.buildHandler()
	container.startJobs(ctx)

	log.Printf("Infra iniciada: DB na porta %s e Kafka no tópico %s",
		container.Config.DBPort, container.Config.QueueTopic)

	return container
}

func (container *Container) buildDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		container.Config.DBHost, container.Config.DBPort, container.Config.DBUser,
		container.Config.DBPassword, container.Config.DBName, container.Config.DBSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	container.DB = db
}

func (container *Container) buildMigration() {
	sourceDriver, err := iofs.New(microum.MigrationsFS, "migrations")
	if err != nil {
		panic(fmt.Errorf("could not create source driver: %w", err))
	}

	dbDriver, err := postgres.WithInstance(container.DB, &postgres.Config{})
	if err != nil {
		panic(fmt.Errorf("could not create database driver: %w", err))
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"postgres",
		dbDriver,
	)
	if err != nil {
		panic(fmt.Errorf("could not create migrate instance: %w", err))
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}

	log.Println("✅ Migrations executed successfully!")
}

func (container *Container) buildKafka() {
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(container.Config.QueueBroker),
		Topic:    container.Config.QueueTopic,
		Balancer: &kafka.LeastBytes{}, // Distribui mensagens entre partições
	}

	container.Kafka = queue.NewKafkaAdapter(kafkaWriter)
}

func (container *Container) buildPrometheus() {
	container.RegCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "m1_processed_registrations_total",
		Help: "The total number of successfully processed registrations",
	})

	container.ErrCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "m1_queue_delivery_errors_total",
		Help: "The total number of errors when delivering messages to queue",
	})
}

func (container *Container) buildRepository() {
	container.Repository = process.NewRepository(container.DB)
}

func (container *Container) buildService() {
	container.Service = process.NewService(container.Repository, container.RegCounter, container.ErrCounter, container.Kafka)
}

func (container *Container) buildHandler() {
	container.Handler = process.NewHandler(container.Service)
}

func (container *Container) startJobs(ctx context.Context) {
	go jobs.NewOutboxJob(container.Repository, container.Kafka, container.ErrCounter).Start(ctx)
}

func (container *Container) Close() {}
