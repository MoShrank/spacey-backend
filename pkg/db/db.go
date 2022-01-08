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

func (db *Database) QueryDocument(
	dbName string,
	collectionName string,
	filter interface{},
) *mongo.SingleResult {
	return db.GetDB(dbName).Collection(collectionName).FindOne(context.TODO(), filter)
}

func (db *Database) QueryDocuments(
	dbName string,
	collectionName string,
	filter interface{},
) (*mongo.Cursor, error) {
	return db.GetDB(dbName).Collection(collectionName).Find(context.TODO(), filter)
}

func (db *Database) CreateDocument(
	dbName string,
	collectionName string,
	document interface{},
) (*mongo.InsertOneResult, error) {
	return db.GetDB(dbName).Collection(collectionName).InsertOne(context.TODO(), document)
}

func (db *Database) UpdateDocument(
	dbName string,
	collectionName string,
	filter interface{},
	update interface{},
) (*mongo.UpdateResult, error) {
	return db.GetDB(dbName).Collection(collectionName).UpdateOne(context.TODO(), filter, update)
}

func (db *Database) DeleteDocument(
	dbName string,
	collectionName string,
	filter interface{},
) (*mongo.DeleteResult, error) {
	return db.GetDB(dbName).Collection(collectionName).DeleteOne(context.TODO(), filter)
}
