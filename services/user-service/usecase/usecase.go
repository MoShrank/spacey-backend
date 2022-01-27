package usecase

import (
	"time"

	mapper "github.com/PeteProgrammer/go-automapper"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/user-service/entity"
)

type UserUsecase struct {
	logger    logger.LoggerInterface
	jwt       auth.JWTInterface
	userStore entity.UserStoreInterface
}

func NewUserUseCase(
	loggerObj logger.LoggerInterface,
	jwtObj auth.JWTInterface,
	userStore entity.UserStoreInterface,
) entity.UserUsecaseInterface {
	return &UserUsecase{
		logger:    loggerObj,
		jwt:       jwtObj,
		userStore: userStore,
	}
}

func (u *UserUsecase) CreateUser(user interface{}) (*entity.UserResponseModel, error) {
	var dbUser entity.User
	mapper.MapLoose(user, &dbUser)

	hashedPassword, err := u.jwt.HashPassword(dbUser.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	dbUser.Password = hashedPassword
	dbUser.CreatedAtTs = &now
	dbUser.UpdatedAtTs = &now
	dbUser.DeletedAtTs = nil

	id, err := u.userStore.SaveUser(&dbUser)
	if err != nil {
		return nil, entity.ErrEmailAlreadyExists
	}

	dbUser.ID = id

	var respUser entity.UserResponseModel
	mapper.MapLoose(dbUser, &respUser)

	return &respUser, nil
}

func (u *UserUsecase) Login(email, password string) (*entity.UserResponseModel, error) {
	dbUser, err := u.userStore.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if ok, _ := u.jwt.CheckPasswordHash(password, dbUser.Password); !ok {
		return nil, err
	}

	token, err := u.jwt.CreateJWTWithClaims(dbUser.ID)
	if err != nil {
		return nil, err
	}

	var respUser entity.UserResponseModel
	mapper.MapLoose(dbUser, &respUser)
	respUser.Token = token

	return &respUser, nil
}

func (u *UserUsecase) GetUserByID(id string) (*entity.UserResponseModel, error) {
	dbUser, err := u.userStore.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	var respUser entity.UserResponseModel
	mapper.MapLoose(dbUser, &respUser)

	return &respUser, nil
}
