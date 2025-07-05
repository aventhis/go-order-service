package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Kafka    KafkaConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port int
	Env  string
}

type KafkaConfig struct {
	Brokers string
	Topic    string
    GroupID  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found; using environment variables instead")
	}

	appPort, err := strconv.Atoi(getEnv("APP_PORT", "8081"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %v", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432")) // порт БД
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %v", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: appPort,
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "db"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "order_app"),
			Password: getEnv("DB_PASSWORD", "secret"),
			Name:     getEnv("DB_NAME", "order_service"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Kafka: KafkaConfig{
			Brokers: getEnv("KAFKA_BROKERS", "kafka:9092"),
			Topic:   getEnv("KAFKA_TOPIC", "orders"),
    		GroupID: getEnv("KAFKA_GROUP_ID", "order_service"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
