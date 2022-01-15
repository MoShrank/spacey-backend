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
	GraylogConnection      string
	DBName                 string
}

type ConfigInterface interface {
	GetPort() string
	GetMongoDBConnection() string
	GetLogLevel() string
	GetJWTSecret() string
	GetUserSeviceDBNAME() string
	GetFlashcardServiceDBName() string
	GetGrayLogConnection() string
	GetDBName() string
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
		GraylogConnection:      loadEnv("GRAYLOG_CONNECTION", "localhost://localhost:12201"),
		AuthSecretKey:          loadEnvWithoutDefault("AUTH_SECRET_KEY"),
		UserServiceDBName:      loadEnvWithoutDefault("USER_SERVICE_DB_NAME"),
		FlashcardServiceDBName: loadEnvWithoutDefault("FLASHCARD_SERVICE_DB_NAME"),
		DBName:                 loadEnvWithoutDefault("DB_NAME"),
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

func (c *Config) GetUserSeviceDBNAME() string {
	return c.UserServiceDBName
}

func (c *Config) GetFlashcardServiceDBName() string {
	return c.FlashcardServiceDBName
}

func (c *Config) GetGrayLogConnection() string {
	return c.GraylogConnection
}

func (c *Config) GetDBName() string {
	return c.DBName
}
