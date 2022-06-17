package usecase

import (
	"fmt"
	"strings"
	"time"

	mapper "github.com/PeteProgrammer/go-automapper"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/user-service/entity"
	"github.com/moshrank/spacey-backend/services/user-service/external"
)

type UserUsecase struct {
	logger      logger.LoggerInterface
	jwt         auth.JWTInterface
	userStore   entity.UserStoreInterface
	emailSender external.EmailSenderInterface
	cfg         config.ConfigInterface
}

func NewUserUseCase(
	loggerObj logger.LoggerInterface,
	jwtObj auth.JWTInterface,
	userStore entity.UserStoreInterface,
	emailSender external.EmailSenderInterface,
	cfg config.ConfigInterface,
) entity.UserUsecaseInterface {
	return &UserUsecase{
		logger:      loggerObj,
		jwt:         jwtObj,
		userStore:   userStore,
		emailSender: emailSender,
		cfg:         cfg,
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
	dbUser.Email = strings.ToLower(dbUser.Email)
	dbUser.CreatedAtTs = &now
	dbUser.UpdatedAtTs = &now
	dbUser.DeletedAtTs = nil
	dbUser.BetaUser = false

	id, err := u.userStore.SaveUser(&dbUser)
	if err != nil {
		return nil, entity.ErrEmailAlreadyExists
	}

	dbUser.ID = id

	var respUser entity.UserResponseModel
	mapper.MapLoose(dbUser, &respUser)

	return &respUser, nil
}

func (u *UserUsecase) CreateToken(id string, isBeta, emailVerified bool) (string, error) {
	token, err := u.jwt.CreateJWTWithClaims(id, map[string]interface{}{
		"IsBeta":         isBeta,
		"EmailValidated": emailVerified,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) Login(email, password string) (*entity.UserResponseModel, error) {
	email = strings.ToLower(email)

	dbUser, err := u.userStore.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if ok, err := u.jwt.CheckPasswordHash(password, dbUser.Password); !ok {
		return nil, err
	}

	token, err := u.CreateToken(dbUser.ID, dbUser.BetaUser, dbUser.EmailValidated)
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

func (u *UserUsecase) VerifyEmail(userID, token string) (string, error) {
	claims, err := u.jwt.ValidateJWT(token, []string{})

	if err != nil {
		return "", err
	}

	id := claims["Id"].(string)

	if id != userID {
		return "", fmt.Errorf("userID mismatch")
	}

	dbUser, err := u.userStore.GetUserByID(id)
	if err != nil {
		return "", err
	}

	dbUser.EmailValidated = true
	err = u.userStore.VerifyEmail(id)
	if err != nil {
		return "", err
	}

	authToken, err := u.CreateToken(id, dbUser.BetaUser, dbUser.EmailValidated)

	return authToken, err
}

func (u *UserUsecase) SendVerificationEmail(id string) error {
	dbUser, err := u.userStore.GetUserByID(id)
	if err != nil {
		return err
	}

	if dbUser.EmailValidated {
		return fmt.Errorf("email already verified")
	}

	token, err := u.jwt.CreateJWTWithClaims(id)
	if err != nil {
		return err
	}

	domain := u.cfg.GetDomain()

	link := fmt.Sprintf("https://%s/verifying?token=%s", domain, token)

	return u.emailSender.SendEmail(dbUser.Email, link)
}
