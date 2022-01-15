package store

import (
	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/services/user-service/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const USER_COLLECTION = "users"

type Store struct {
	db db.DatabaseInterface
}

func NewStore(db db.DatabaseInterface) entity.UserStoreInterface {
	return &Store{
		db: db,
	}
}

func (s *Store) SaveUser(user *entity.User) (string, error) {
	res, err := s.db.CreateDocument(USER_COLLECTION, user)

	if err != nil {
		return "", err
	}

	id, err := res.InsertedID.(primitive.ObjectID).MarshalJSON()

	return string(id[:]), err
}

func (s *Store) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	res := s.db.QueryDocument(USER_COLLECTION, map[string]interface{}{"email": email})
	err := res.Decode(&user)

	return &user, err
}
