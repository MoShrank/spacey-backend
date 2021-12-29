package usecase

import (
	"testing"

	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	var usecase = NewUserUseCase(logger.NewLogger(""), SecretKey{SecretKey: []byte("secret")})
	password := "password"
	hashedPassword, err := usecase.HashPassword(password)

	assert.NotEqual(t, password, hashedPassword)
	assert.Equal(t, nil, err)
	assert.Greater(t, len(hashedPassword), 0)
}

func TestCheckPasswordHash(t *testing.T) {
	var usecase = NewUserUseCase(logger.NewLogger(""), SecretKey{SecretKey: []byte("secret")})
	password := "password"
	hash := "password"

	assert.False(t, usecase.CheckPasswordHash(password, hash))
}

func TestValidateJWT(t *testing.T) {
	var usecase = NewUserUseCase(logger.NewLogger(""), SecretKey{SecretKey: []byte("secret")})
	jwtString := "test"

	ok, err := usecase.ValidateJWT(jwtString)

	assert.Equal(t, false, ok)
	assert.NotEqual(t, nil, err)
}

func TestCreateJWTWithClaims(t *testing.T) {
	var usecase = NewUserUseCase(logger.NewLogger(""), SecretKey{SecretKey: []byte("secret")})

	jwtString, err := usecase.CreateJWTWithClaims("test")

	assert.Equal(t, nil, err)
	assert.Greater(t, len(jwtString), 0)
}
