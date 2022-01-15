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
	DB     *mongo.Database
}

type DatabaseInterface interface {
	GetDB(string) *mongo.Database
	connect(string) (*mongo.Client, error)
	QueryDocument(string, interface{}) *mongo.SingleResult
	QueryDocuments(string, interface{}) (*mongo.Cursor, error)
	CreateDocument(string, interface{}) (*mongo.InsertOneResult, error)
	UpdateDocument(string, interface{}, interface{}) (*mongo.UpdateResult, error)
	DeleteDocument(string, interface{}) (*mongo.DeleteResult, error)
}

func NewDB(connectionString, dbName string, logger logger.LoggerInterface) DatabaseInterface {
	db := &Database{
		client: nil,
		logger: logger,
		DB:     nil,
	}

	client, err := db.connect(connectionString)

	if err != nil {
		logger.Fatal("Could not connect to Database:", err)
	}

	db.client = client
	db.DB = db.client.Database(dbName)

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
	collectionName string,
	filter interface{},
) *mongo.SingleResult {
	return db.DB.Collection(collectionName).FindOne(context.TODO(), filter)
}

func (db *Database) QueryDocuments(
	collectionName string,
	filter interface{},
) (*mongo.Cursor, error) {
	return db.DB.Collection(collectionName).Find(context.TODO(), filter)
}

func (db *Database) CreateDocument(
	collectionName string,
	document interface{},
) (*mongo.InsertOneResult, error) {
	return db.DB.Collection(collectionName).InsertOne(context.TODO(), document)
}

func (db *Database) UpdateDocument(
	collectionName string,
	filter interface{},
	update interface{},
) (*mongo.UpdateResult, error) {
	return db.DB.Collection(collectionName).UpdateOne(context.TODO(), filter, update)
}

func (db *Database) DeleteDocument(
	collectionName string,
	filter interface{},
) (*mongo.DeleteResult, error) {
	return db.DB.Collection(collectionName).DeleteOne(context.TODO(), filter)
}
