package db

import (
	"context"

	"github.com/moshrank/spacey-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	db     *mongo.Database
	logger logger.LoggerInterface
}

type DatabaseInterface interface {
	GetDB() *mongo.Database
}

func NewDB(connectionString string, logger logger.LoggerInterface) DatabaseInterface {
	database, err := connect(connectionString)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Database Connection Established!")

	return &Database{
		db:     database,
		logger: logger,
	}
}

func connect(connectionString string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	database := client.Database("users")

	return database, nil
}

func (db *Database) GetDB() *mongo.Database {
	return db.db
}
