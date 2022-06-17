package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup(exp int) JWT {
	return JWT{
		secretKey:    []byte("secret"),
		expireOffset: time.Second * time.Duration(exp),
	}
}

func TestCreateJWT(t *testing.T) {
	jwtObj := setup(60)
	token, err := jwtObj.CreateJWTWithClaims("1")
	assert.NotEmpty(t, token)
	assert.NoError(t, err)
}

func TestValidateValidJWT(t *testing.T) {
	jwtObj := setup(60)
	token, _ := jwtObj.CreateJWTWithClaims("1")
	claims, err := jwtObj.ValidateJWT(token, []string{})
	assert.NoError(t, err)
	assert.NotNil(t, claims)
}

func TestValidateInvalidJWT(t *testing.T) {
	// set expire date to 0 so token is expired
	jwtObj := setup(0)
	token, _ := jwtObj.CreateJWTWithClaims("1")
	// sleep one second to invalidate token
	time.Sleep(time.Second * 2)
	claims, err := jwtObj.ValidateJWT(token, []string{})
	assert.Nil(t, claims)
	assert.Error(t, err)
}

func TestHashPassword(t *testing.T) {
	jwtObj := setup(0)
	hash, err := jwtObj.HashPassword("password")

	assert.NotEmpty(t, hash)
	assert.NoError(t, err)
}

func TestCheckPasswordHash(t *testing.T) {
	jwtObj := setup(0)
	hash, _ := jwtObj.HashPassword("password")
	ok, err := jwtObj.CheckPasswordHash("password", hash)

	assert.True(t, ok)
	assert.NoError(t, err)
}
