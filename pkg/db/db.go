package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	db *mongo.Database
}

type DatabaseInterface interface {
}

func NewDB(connectionString string) DatabaseInterface {
	database, err := connect(connectionString)

	if err != nil {
		log.Panic(err)
	}

	log.Print("Database Connection Established!")

	return &Database{
		db: database,
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
