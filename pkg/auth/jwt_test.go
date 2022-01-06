package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup(exp int) JWT {
	return JWT{
		secretKey:    []byte("secret"),
		expireOffset: 60 * 60 * 24,
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
	ok, err := jwtObj.ValidateJWT(token)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestValidateInvalidJWT(t *testing.T) {
	// set expire date to 0 so token is expired
	jwtObj := setup(0)
	token, _ := jwtObj.CreateJWTWithClaims("1")
	// sleep one second to invalidate token
	time.Sleep(time.Second)
	ok, err := jwtObj.ValidateJWT(token)

	assert.False(t, ok)
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

func CheckExtractClaimsValidClaims(t *testing.T) {
	jwtObj := setup(0)

	token, _ := jwtObj.CreateJWTWithClaims("1")

	claims, ok := jwtObj.ExtractClaims(token)

	assert.True(t, ok)
	assert.Equal(t, "1", claims["sub"])
}

func CheckExtractClaimsInvalidClaims(t *testing.T) {
	jwtObj := setup(0)

	token := "invalid"

	claims, ok := jwtObj.ExtractClaims(token)

	assert.False(t, ok)
	assert.Equal(t, nil, claims)
}
