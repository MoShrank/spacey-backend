package db

import (
	"context"
	"embed"

	_ "embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/pkg/errors"
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

//go:embed migrations/*
var migrationFiles embed.FS

func NewDB(cfg config.ConfigInterface, logger logger.LoggerInterface) DatabaseInterface {
	db := &Database{
		client: nil,
		logger: logger,
		DB:     nil,
	}

	db.logger.Info("Running migration...")
	err := db.runMigration(
		cfg.GetMongoDBConnection(),
	)
	if err != nil {
		logger.Error("Could not run migrations: ", err)
	}

	client, err := db.connect(cfg.GetMongoDBConnection())

	if err != nil {
		logger.Fatal("Could not connect to Database: ", err)
	}

	db.client = client
	db.DB = db.client.Database(cfg.GetDBName())

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

func (db *Database) runMigration(connString string) error {
	driver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return errors.Wrap(err, "could not create migration driver")
	}
	m, err := migrate.NewWithSourceInstance("iofs", driver, connString)
	if err != nil {
		return err
	}

	err = m.Up()

	return err
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
	return db.DB.Collection(collectionName).
		UpdateOne(context.TODO(), filter, update)
}

func (db *Database) DeleteDocument(
	collectionName string,
	filter interface{},
) (*mongo.DeleteResult, error) {
	return db.DB.Collection(collectionName).DeleteOne(context.TODO(), filter)
}
