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
	connect(string) (*mongo.Client, error)
}

func NewDB(connectionString string, logger logger.LoggerInterface) DatabaseInterface {
	db := &Database{
		client: nil,
		logger: logger,
	}

	client, err := db.connect(connectionString)

	if err != nil {
		logger.Fatal("Could not connect to Database:", err)
	}

	db.client = client
	logger.Info("Database Connection Established!")

	return db
}

func (db *Database) connect(connectionString string) (*mongo.Client, error) {
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
