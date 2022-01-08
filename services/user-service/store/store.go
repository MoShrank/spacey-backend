package store

import (
	"context"

	"github.com/moshrank/spacey-backend/services/user-service/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db *mongo.Database
}

func NewStore(db *mongo.Database) entity.UserStoreInterface {
	return &Store{
		db: db,
	}
}

func (db *Store) SaveUser(user *entity.User) (string, error) {
	userCollection := db.db.Collection("users")

	insertionResult, err := userCollection.InsertOne(context.TODO(), user)

	return insertionResult.InsertedID.(string), err
}

func (db *Store) GetUserByEmail(email string) (*entity.User, error) {
	userCollection := db.db.Collection("users")

	var user entity.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	return &user, err
}
