package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	MongoDBConnection   string
	LogLevel            string
	AuthSecretKey       string
	GraylogConnection   string
	DBName              string
	UserServiceHostName string
	deckServiceHostName string
}

type ConfigInterface interface {
	GetPort() string
	GetMongoDBConnection() string
	GetLogLevel() string
	GetJWTSecret() string
	GetGrayLogConnection() string
	GetDBName() string
	GetUserServiceHostName() string
	GetDeckServiceHostName() string
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
		Port:                loadEnv("PORT", "8080"),
		MongoDBConnection:   loadEnv("MONGO_DB_CONNECTION", "mongodb://127.0.0.1:27017"),
		LogLevel:            loadEnv("LOG_LEVEL", "debug"),
		GraylogConnection:   loadEnv("GRAYLOG_CONNECTION", "localhost://localhost:12201"),
		AuthSecretKey:       loadEnvWithoutDefault("AUTH_SECRET_KEY"),
		DBName:              loadEnvWithoutDefault("DB_NAME"),
		UserServiceHostName: loadEnv("USER_SERVICE_HOST_NAME", "user-service"),
		deckServiceHostName: loadEnv("DECK_SERVICE_HOST_NAME", "deck-service"),
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

func (c *Config) GetJWTSecret() string {
	return c.AuthSecretKey
}

func (c *Config) GetGrayLogConnection() string {
	return c.GraylogConnection
}

func (c *Config) GetDBName() string {
	return c.DBName
}

func (c *Config) GetUserServiceHostName() string {
	return c.UserServiceHostName
}

func (c *Config) GetDeckServiceHostName() string {
	return c.deckServiceHostName
}
