package store

import (
	"github.com/moshrank/spacey-backend/pkg/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const USER_COLLECTION = "user"

type User struct {
	ID             string `bson:"_id,omitempty"`
	BetaUser       bool   `bson:"betaUser"`
	EmailValidated bool   `bson:"emailValidated"`
}

type StoreInterface interface {
	GetUserByID(string) (*User, error)
}

type Store struct {
	db db.DatabaseInterface
}

func NewStore(db db.DatabaseInterface) StoreInterface {
	return &Store{db: db}
}

func (s *Store) GetUserByID(id string) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.db.QueryDocument(USER_COLLECTION, map[string]interface{}{"_id": objectID})
	if err != nil {
		return nil, err
	}

	var user User
	err = res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
