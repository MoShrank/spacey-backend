package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                          string
	MongoDBConnection             string
	LogLevel                      string
	AuthSecretKey                 string
	GraylogConnection             string
	DBName                        string
	UserServiceHostName           string
	DeckServiceHostName           string
	LearningServiceHostName       string
	CardGenerationServiceHostName string
	Domain                        string
	MaxAgeAuth                    int
	MigrationFilePath             string
	MailGunAPIKey                 string
	Environment                   string
}

type ConfigInterface interface {
	GetPort() string
	GetMongoDBConnection() string
	GetLogLevel() string
	GetJWTSecret() string
	GetDBName() string
	GetUserServiceHostName() string
	GetDeckServiceHostName() string
	GetLearningServiceHostName() string
	GetDomain() string
	GetMaxAgeAuth() int
	GetCardGenerationServiceHostName() string
	GetMailGunAPIKey() string
	GetEnv() string
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
	envFile := os.Getenv("ENV_FILE_PATH")

	if envFile == "" {
		envFile = ".env"
	}

	err := godotenv.Load(envFile)

	if err != nil {
		fmt.Println("Warning: could not find .env file under: ", envFile)
	}

	maxAgeAuth, err := strconv.Atoi(loadEnv("MAX_AGE_AUTH", "604800"))

	if err != nil {
		panic(err)
	}

	return &Config{
		Port: loadEnv("PORT", "8080"),
		MongoDBConnection: loadEnv(
			"MONGO_DB_CONNECTION",
			"mongodb://mongodb:27017/spacey",
		),
		LogLevel:                loadEnv("LOG_LEVEL", "info"),
		AuthSecretKey:           loadEnv("AUTH_SECRET_KEY", "test_secret_key"),
		DBName:                  loadEnv("DB_NAME", "spacey"),
		UserServiceHostName:     loadEnv("USER_SERVICE_HOST_NAME", "user-service"),
		DeckServiceHostName:     loadEnv("DECK_SERVICE_HOST_NAME", "deck-management-service"),
		Domain:                  loadEnv("DOMAIN", "localhost"),
		MaxAgeAuth:              maxAgeAuth,
		LearningServiceHostName: loadEnv("LEARNING_SERVICE_HOST_NAME", "learning-service"),
		CardGenerationServiceHostName: loadEnv(
			"CARD_GENERATION_SERVICE_HOST_NAME",
			"card-generation-service",
		),
		MailGunAPIKey: loadEnv("MAIL_GUN_API_KEY", ""),
		Environment:   loadEnv("ENVIRONMENT", "dev"),
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

func (c *Config) GetLearningServiceHostName() string {
	return c.LearningServiceHostName
}

func (c *Config) GetCardGenerationServiceHostName() string {
	return c.CardGenerationServiceHostName
}

func (c *Config) GetMailGunAPIKey() string {
	return c.MailGunAPIKey
}

func (c *Config) GetEnv() string {
	return c.Environment
}
