package infra

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment  string
	ServerPort   string
	ServerName   string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	QueueBroker  string
	QueueTopic   string
	OtelEndpoint string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		Environment:  os.Getenv("ENVIRONMENT"),
		ServerPort:   os.Getenv("SERVER_PORT"),
		ServerName:   os.Getenv("SERVER_NAME"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		DBSSLMode:    os.Getenv("DB_SSLMODE"),
		QueueBroker:  os.Getenv("KAFKA_BROKER"),
		QueueTopic:   os.Getenv("KAFKA_TOPIC"),
		OtelEndpoint: os.Getenv("OTEL_ENDPOINT"),
	}
}
