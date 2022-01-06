package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type JWT struct {
	secretKey    []byte
	expireOffset int64
}

type JWTInterface interface {
	ExtractClaims(tokenStr string) (jwt.MapClaims, bool)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) (bool, error)
	ValidateJWT(tokenString string) (bool, error)
	CreateJWTWithClaims(userID string) (string, error)
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		secretKey:    []byte(secretKey),
		expireOffset: 60 * 60 * 24,
	}
}

func (j *JWT) ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return nil, false
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

func (j *JWT) ValidateJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["sub"] != nil &&
		claims.VerifyExpiresAt(j.expireOffset, true) {
		return true, nil
	} else {
		return false, err
	}
}

func (j *JWT) CreateJWTWithClaims(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(j.expireOffset)).Unix(),
	})

	tokenString, err := token.SignedString(j.secretKey)
	return tokenString, err
}
