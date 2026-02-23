package infra

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment   string
	ServerName    string
	ServerPort    string
	KafkaAddr     string
	TopicInput    string
	TopicOutput   string
	ConsumerGroup string
	MongoURI      string
	MongoDatabase string
	OtelEndpoint  string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		Environment:   os.Getenv("ENVIRONMENT"),
		ServerName:    "",
		ServerPort:    "",
		KafkaAddr:     "",
		TopicInput:    "",
		TopicOutput:   "",
		ConsumerGroup: "",
		MongoURI:      "",
		MongoDatabase: "",
		OtelEndpoint:  "",
	}
}
