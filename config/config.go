package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	MongoDBConnection      string
	LogLevel               string
	AuthSecretKey          string
	UserServiceDBName      string
	FlashcardServiceDBName string
}

type ConfigInterface interface {
	GetPort() string
	GetMongoDBConnection() string
	GetLogLevel() string
	GetSecretKey() string
	GetUserSeviceDBNAME() string
	GetFlashcardServiceDBName() string
}

func loadEnvWithoutDefault(key string) string {
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
	godotenv.Load()

	return &Config{
		Port:                   loadEnv("PORT", "8080"),
		MongoDBConnection:      loadEnv("MONGO_DB_CONNECTION", "mongodb://mongodb:27017"),
		LogLevel:               loadEnv("LOG_LEVEL", "debug"),
		AuthSecretKey:          loadEnvWithoutDefault("AUTH_SECRET_KEY"),
		UserServiceDBName:      loadEnvWithoutDefault("USER_SERVICE_DB_NAME"),
		FlashcardServiceDBName: loadEnvWithoutDefault("FLASHCARD_SERVICE_DB_NAME"),
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

func (c *Config) GetUserSeviceDBNAME() string {
	return c.UserServiceDBName
}

func (c *Config) GetFlashcardServiceDBName() string {
	return c.FlashcardServiceDBName
}
