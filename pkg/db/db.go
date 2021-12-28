package db

import (
	"context"

	"github.com/moshrank/spacey-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client *mongo.Client
	logger logger.LoggerInterface
}

type DatabaseInterface interface {
	GetDB(string) *mongo.Database
}

func NewDB(connectionString string, logger logger.LoggerInterface) DatabaseInterface {
	client, err := connect(connectionString)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Database Connection Established!")

	return &Database{
		client: client,
		logger: logger,
	}
}

func connect(connectionString string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (db *Database) GetDB(dbName string) *mongo.Database {
	return db.client.Database(dbName)

}
