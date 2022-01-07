package store

import (
	"context"

	"github.com/moshrank/spacey-backend/services/user-service/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db *mongo.Database
}

type StoreInterface interface {
	SaveUser(user *entity.User) error
	GetPassword(email string) (string, error)
	GetUserByEmail(email string) (*entity.User, error)
}

func NewStore(db *mongo.Database) StoreInterface {
	return &Store{
		db: db,
	}
}

func (db *Store) SaveUser(user *entity.User) error {
	userCollection := db.db.Collection("users")

	_, err := userCollection.InsertOne(context.TODO(), user)

	return err
}

func (db *Store) GetPassword(email string) (string, error) {
	userCollection := db.db.Collection("users")
	var user entity.User
	err := userCollection.FindOne(context.TODO(), map[string]interface{}{"email": email}).
		Decode(&user)

	return user.Password, err
}

func (db *Store) GetUserByEmail(email string) (*entity.User, error) {
	userCollection := db.db.Collection("users")
	var user entity.User
	err := userCollection.FindOne(context.TODO(), map[string]interface{}{"email": email}).
		Decode(&user)

	return &user, err
}
