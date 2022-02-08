package store

import (
	"fmt"

	"github.com/moshrank/spacey-backend/pkg/db"
	"github.com/moshrank/spacey-backend/services/user-service/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const USER_COLLECTION = "user"

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

	id, _ := res.InsertedID.(primitive.ObjectID)

	return id.Hex(), err
}

func (s *Store) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	res := s.db.QueryDocument(USER_COLLECTION, map[string]interface{}{"email": email})
	err := res.Decode(&user)

	return &user, err
}

func (s *Store) GetUserByID(id string) (*entity.User, error) {
	var user entity.User

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid id. cannot convert id into object id: %s error: %w",
			id,
			err,
		)
	}

	res := s.db.QueryDocument(USER_COLLECTION, map[string]interface{}{"_id": objectId})
	err = res.Decode(&user)

	return &user, err
}
