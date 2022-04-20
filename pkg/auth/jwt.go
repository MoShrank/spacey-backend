package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/moshrank/spacey-backend/config"
	"golang.org/x/crypto/bcrypt"
)

type JWT struct {
	secretKey    []byte
	expireOffset time.Duration
}

type JWTInterface interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) (bool, error)
	ValidateJWT(tokenString string) (jwt.MapClaims, error)
	CreateJWTWithClaims(userID string, isBeta bool) (string, error)
}

func NewJWT(cfg config.ConfigInterface) JWTInterface {
	secretKey := []byte(cfg.GetJWTSecret())

	return &JWT{
		secretKey:    secretKey,
		expireOffset: time.Second * time.Duration(cfg.GetMaxAgeAuth()),
	}
}

func (j *JWT) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (j *JWT) CheckPasswordHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, err
}

func (j *JWT) ValidateJWT(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			if !token.Claims.(jwt.MapClaims).VerifyExpiresAt(time.Now().Unix(), true) {
				return false, fmt.Errorf("token expired")
			}

			return j.secretKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error while parsing token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error while asserting standard claims: %w", err)
	}

	if claims["Id"] == "" {
		return nil, fmt.Errorf("invalid token, cannot find user id")
	}

	return claims, nil
}

func (j *JWT) CreateJWTWithClaims(userID string, isBeta bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":     userID,
		"exp":    time.Now().Add(j.expireOffset).Unix(),
		"IsBeta": isBeta,
	})

	tokenString, err := token.SignedString(j.secretKey)
	return tokenString, err
}
