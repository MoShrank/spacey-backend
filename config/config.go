package config

import (
	"errors"
	"os"
)

type Config struct {
	Port              string
	MongoDBConnection string
	LogLevel          string
	AuthSecretKey     string
}

type ConfigInterface interface {
	GetPort() string
	GetMongoDBConnection() string
	GetLogLevel() string
	GetSecretKey() string
}

func loadEnvWithDefault(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(errors.New("Missing env variable: " + key))
	}

	return value
}

func loadEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}

	return value
}

func NewConfig() (ConfigInterface, error) {
	return &Config{
		Port:              loadEnv("PORT", "8080"),
		MongoDBConnection: loadEnv("MONGO_DB_CONNECTION", "mongodb://localhost:27017"),
		LogLevel:          loadEnv("LOG_LEVEL", "debug"),
		AuthSecretKey:     loadEnvWithDefault("AUTH_SECRET_KEY"),
	}, nil
}

func (c *Config) GetPort() string {
	return c.Port
}

func (c *Config) GetMongoDBConnection() string {
	return c.MongoDBConnection
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

func (c *Config) GetSecretKey() string {
	return c.AuthSecretKey
}
