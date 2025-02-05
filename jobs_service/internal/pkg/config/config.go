package config

import (
	"os"
	"strings"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SslMode  string
	}
	Mongodb struct {
		Host     string
		Port     string
		Database string
	}

	OTLPCollector struct {
		Host string
		Port string
	}

	Kafka struct {
		Address []string
		Topic   struct {
			Client string
		}
	}
}

func New() *Config {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.RPCPort = getEnv("RPC_PORT", ":4040")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "db")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "123")
	config.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "ekzamen5_db")

	// mongodb configuration
	config.Mongodb.Host = getEnv("POSTGRES_HOST", "mongodb")
	config.Mongodb.Port = getEnv("POSTGRES_PORT", "27017")
	config.Mongodb.Database = getEnv("POSTGRES_DATABASE", "ekzamen5_db")

	// otlp collector configuration
	config.OTLPCollector.Host = getEnv("OTLP_COLLECTOR_HOST", "otel-collector")
	config.OTLPCollector.Port = getEnv("OTLP_COLLECTOR_PORT", ":4317")

	// kafka configuration
	config.Kafka.Address = strings.Split(getEnv("KAFKA_ADDRESS", "kafka:29092"), ",")
	config.Kafka.Topic.Client = getEnv("KAFKA_TOPIC_CLIENT_CREATE", "jobs.created")

	return &config
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
