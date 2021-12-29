package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/pkg/logger"
	"github.com/moshrank/spacey-backend/services/user-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DBMock struct {
	mock.Mock
}

func (m *DBMock) GetPassword(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

func (m *DBMock) SaveUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *DBMock) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

type UserUsecaseMock struct {
	mock.Mock
}

func (us *UserUsecaseMock) HashPassword(password string) (string, error) {
	args := us.Called(password)
	return args.String(0), args.Error(1)
}

func (us *UserUsecaseMock) CheckPasswordHash(password string, hash string) bool {
	args := us.Called(password, hash)
	return args.Bool(0)
}

func (us *UserUsecaseMock) ValidateJWT(tokenString string) (bool, error) {
	args := us.Called(tokenString)
	return args.Bool(0), args.Error(1)
}

func (us *UserUsecaseMock) CreateJWTWithClaims(userID string) (string, error) {
	args := us.Called(userID)
	return args.String(0), args.Error(1)
}

func TestCreateUserHandler(t *testing.T) {

	tests := []struct {
		TestName   string
		body       string
		statusCode int
	}{
		{
			"Valid User",
			"{\"name\": \"moritz\", \"email\": \"moritz.e50@gmail.com\", \"password\": \"test_password\"}",
			201,
		},
		{
			"Empty Body",
			"",
			400,
		},
		{
			"Empty Email",
			"{\"name\": \"moritz\", \"email\": \", \"password\": \"test_password\"}",
			400,
		},
		{
			"Invalid Email",
			"{\"name\": \"moritz\", \"email\": \"moritz.e50gmail.com\", \"password\": \"test_password\"}",
			400,
		},
		{
			"Empty Password",
			"{\"name\": \"moritz\", \"email\": \"moritz.e50@gmail.com\", \"password\": \"\"}",
			400,
		},
		{
			"Short Password",
			"{\"name\": \"moritz\", \"email\": \"moritz.e50@gmail.com\", \"password\": \"test\"}",
			400,
		},
		{
			"Empty Name",
			"{\"name\": \"\", \"email\": \"moritz.e50@gmail.com\", \"password\": \"test_password\"}",
			400,
		},
	}

	dbMock := new(DBMock)
	userUsecaseMock := new(UserUsecaseMock)

	var handler = NewHandler(
		dbMock,
		logger.NewLogger(""),
		userUsecaseMock,
	)

	dbMock.On("SaveUser", mock.Anything).Return(nil)
	userUsecaseMock.On("HashPassword", mock.Anything).Return("", nil)

	for _, test := range tests {

		t.Run(test.TestName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/user", bytes.NewBuffer([]byte(test.body)))
			handler.CreateUser(c)

			assert.Equal(t, test.statusCode, c.Writer.Status())
		})
	}
}

func TestLogin(t *testing.T) {

	tests := []struct {
		TestName   string
		body       string
		statusCode int
	}{
		{
			"Valid User",
			"{\"email\": \"moritz.e50@gmail.com\", \"password\": \"test_password\"}",
			200,
		},
		{
			"Empty Body",
			"",
			400,
		},
		{
			"Empty Email",
			"{\"email\": \", \"password\": \"test_password\"}",
			400,
		},
		{
			"Empty Password",
			"{\"email\": \"moritz.e50@gmail.com\", \"password\": \"\"}",
			400,
		},
	}

	dbMock := new(DBMock)
	userUsecaseMock := new(UserUsecaseMock)

	var handler = NewHandler(
		dbMock,
		logger.NewLogger(""),
		userUsecaseMock,
	)

	dbMock.On("GetUserByEmail", mock.Anything).Return(&models.User{}, nil)
	userUsecaseMock.On("CheckPasswordHash", mock.Anything, mock.Anything).Return(true)
	userUsecaseMock.On("CreateJWTWithClaims", mock.Anything).Return("", nil)

	for _, test := range tests {

		t.Run(test.TestName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/user/login", bytes.NewBuffer([]byte(test.body)))
			handler.Login(c)

			assert.Equal(t, test.statusCode, c.Writer.Status())
		})
	}
}
