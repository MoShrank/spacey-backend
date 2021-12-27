package usecase

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type SecretKey struct {
	SecretKey string
}

type UserUsecase struct {
	logger    logger.LoggerInterface
	secretKey string
}

type UserUsecaseInterface interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	ValidateJWT(tokenString string) (bool, error)
	CreateJWTWithClaims(usesrID string) (string, error)
}

func NewUserUseCase(loggerObj logger.LoggerInterface, secretKey SecretKey) UserUsecaseInterface {
	return &UserUsecase{
		logger:    loggerObj,
		secretKey: secretKey.SecretKey,
	}
}

func (u *UserUsecase) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UserUsecase) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	u.logger.Error(err.Error())
	return err == nil
}

func (u *UserUsecase) ValidateJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return u.secretKey, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	} else {
		return false, err
	}

}

func (u *UserUsecase) CreateJWTWithClaims(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(u.secretKey)

	return tokenString, err
}
