package config

import (
	"errors"
	"os"
	"strconv"

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
	DeckServiceHostName string
	Domain              string
	MaxAgeAuth          int
	MigrationFilePath   string
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
	GetDomain() string
	GetMaxAgeAuth() int
	GetMigrationFilePath() string
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

	maxAgeAuth, err := strconv.Atoi(loadEnv("MAX_AGE_AUTH", "604800"))

	if err != nil {
		panic(err)
	}

	return &Config{
		Port:                loadEnv("PORT", "8080"),
		MongoDBConnection:   loadEnv("MONGO_DB_CONNECTION", "mongodb://127.0.0.1:27017/spacey"),
		LogLevel:            loadEnv("LOG_LEVEL", "info"),
		GraylogConnection:   loadEnv("GRAYLOG_CONNECTION", "localhost://localhost:12201"),
		AuthSecretKey:       loadEnv("AUTH_SECRET_KEY", "secret"),
		DBName:              loadEnv("DB_NAME", "spacey"),
		UserServiceHostName: loadEnv("USER_SERVICE_HOST_NAME", "user-service"),
		DeckServiceHostName: loadEnv("DECK_SERVICE_HOST_NAME", "deck-management-service"),
		Domain:              loadEnv("DOMAIN", "localhost"),
		MaxAgeAuth:          maxAgeAuth,
		MigrationFilePath:   loadEnv("MIGRATION_FILE_PATH", "../../migrations"),
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
	return c.DeckServiceHostName
}

func (c *Config) GetDomain() string {
	return c.Domain
}

func (c *Config) GetMaxAgeAuth() int {
	return c.MaxAgeAuth
}

func (c *Config) GetMigrationFilePath() string {
	return c.MigrationFilePath
}
