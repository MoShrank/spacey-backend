package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/httpconst"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/user-service/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (u *UserUsecaseMock) CreateUser(user interface{}) (*entity.UserResponseModel, error) {
	args := u.Called(user)
	return args.Get(0).(*entity.UserResponseModel), args.Error(1)
}

func (u *UserUsecaseMock) Login(email, password string) (*entity.UserResponseModel, error) {
	args := u.Called(email, password)
	return args.Get(0).(*entity.UserResponseModel), args.Error(1)
}

func (u *UserUsecaseMock) GetUserByID(id string) (*entity.UserResponseModel, error) {
	args := u.Called(id)
	return args.Get(0).(*entity.UserResponseModel), args.Error(1)
}

func getJSONErr(statusCode int) string {
	return fmt.Sprintf("{\"error\": \"%s\", \"message\": \"\"}", httpconst.ErrorMapping[statusCode])
}

type Err struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func TestCreateUserHandlerInvalidUser(t *testing.T) {

	tests := []struct {
		TestName       string
		inpBody        string
		wantBody       string
		wantStatusCode int
	}{
		{
			"Empty Body",
			"",
			getJSONErr(400),
			400,
		},
		{
			"Empty Email",
			"{\"name\": \"moritz\", \"email\": \"\", \"password\": \"test_password\"}",
			getJSONErr(400),
			400,
		},
		{
			"Invalid Email",
			"{\"name\": \"moritz\", \"email\": \"moritz.e50gmail.com\", \"password\": \"test_password\"}",
			getJSONErr(400),
			400,
		},
		{
			"Empty Password",
			"{\"name\": \"moritz\", \"email\": \"moritz.e50@gmail.com\", \"password\": \"\"}",
			getJSONErr(400),
			400,
		},
		{
			"Short Password",
			"{\"name\": \"moritz\", \"email\": \"moritz.e50@gmail.com\", \"password\": \"test\"}",
			getJSONErr(400),
			400,
		},
		{
			"Empty Name",
			"{\"name\": \"\", \"email\": \"moritz.e50@gmail.com\", \"password\": \"test_password\"}",
			getJSONErr(400),
			400,
		},
	}

	usecaseMock := &UserUsecaseMock{}
	conf, _ := config.NewConfig()

	handler := Handler{
		logger:      log.New(),
		userUsecase: usecaseMock,
		validator:   validator.NewValidator(),
		config:      conf,
	}
	usecaseMock.On("CreateUser", mock.Anything).Return(&entity.UserResponseModel{}, nil)
	usecaseMock.On("Login", mock.Anything, mock.Anything).Return(&entity.UserResponseModel{}, nil)

	for _, test := range tests {

		t.Run(test.TestName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/user", bytes.NewBuffer([]byte(test.inpBody)))
			handler.CreateUser(c)

			// set message to empty strint since function does not return custom error msg
			bodyRes := &Err{}
			json.Unmarshal(w.Body.Bytes(), bodyRes)
			bodyRes.Message = ""
			expBody, _ := json.Marshal(bodyRes)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
			assert.JSONEq(t, test.wantBody, string(expBody))
		})
	}
}

func TestCreateUserValidUser(t *testing.T) {

	inpBody := `{"name": "moritz", "email": "moritz.e50@gmail.com", "password": "test_password"}`
	wantBody := `{"data": {"id": "1", "name": "moritz", "email": "moritz.e50@gmail.com", "token": "test_token", "betaUser": false}, "message": "Created"}`
	wantStatusCode := 201

	usecaseMock := &UserUsecaseMock{}
	conf, _ := config.NewConfig()

	handler := Handler{
		logger:      log.New(),
		userUsecase: usecaseMock,
		validator:   validator.NewValidator(),
		config:      conf,
	}
	usecaseMock.On("CreateUser", mock.Anything).Return(&entity.UserResponseModel{
		ID:    "1",
		Name:  "moritz",
		Email: "moritz.e50@gmail.com",
	}, nil)

	usecaseMock.On("Login", mock.Anything, mock.Anything).Return(&entity.UserResponseModel{
		ID:       "1",
		Name:     "moritz",
		Email:    "moritz.e50@gmail.com",
		Token:    "test_token",
		BetaUser: false,
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/user", bytes.NewBuffer([]byte(inpBody)))
	handler.CreateUser(c)

	assert.Equal(t, wantStatusCode, c.Writer.Status())
	assert.JSONEq(t, wantBody, w.Body.String())

}

func TestLoginInvalidUser(t *testing.T) {

	tests := []struct {
		TestName       string
		inpBody        string
		wantBody       string
		wantStatusCode int
	}{
		{
			"Empty Body",
			"",
			getJSONErr(400),
			400,
		},
		{
			"Empty Email",
			"{\"email\": \", \"password\": \"test_password\"}",
			getJSONErr(400),
			400,
		},
		{
			"Empty Password",
			"{\"email\": \"moritz.e50@gmail.com\", \"password\": \"\"}",
			getJSONErr(400),
			400,
		},
	}

	userUsecaseMock := new(UserUsecaseMock)
	conf, _ := config.NewConfig()

	var handler = NewHandler(
		log.New(),
		userUsecaseMock,
		validator.NewValidator(),
		conf,
	)

	userUsecaseMock.On("Login", mock.Anything).Return(&entity.User{}, true)

	for _, test := range tests {

		t.Run(test.TestName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(
				"GET",
				"/user/login",
				bytes.NewBuffer([]byte(test.inpBody)),
			)
			handler.Login(c)

			// set message to empty strint since function does not return custom error msg
			bodyRes := &Err{}
			json.Unmarshal(w.Body.Bytes(), bodyRes)
			bodyRes.Message = ""
			expBody, _ := json.Marshal(bodyRes)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
			assert.JSONEq(t, test.wantBody, string(expBody))
		})
	}
}

func TestLoginValidUser(t *testing.T) {
	inpBody := "{\"email\": \"moritz.e50@gmail.com\", \"password\": \"test_password\"}"
	wantBody := `{"data": 
					{"id": "1", "name": "moritz", "email": "moritz.e50@gmail.com", "token": "test_token", "betaUser": false}, 
					"message": "Success"
				}`
	wantStatusCode := 200

	usecaseMock := &UserUsecaseMock{}
	conf, _ := config.NewConfig()

	handler := Handler{
		logger:      log.New(),
		userUsecase: usecaseMock,
		validator:   validator.NewValidator(),
		config:      conf,
	}
	usecaseMock.On("Login", mock.Anything, mock.Anything).Return(&entity.UserResponseModel{
		ID:    "1",
		Name:  "moritz",
		Email: "moritz.e50@gmail.com",
		Token: "test_token",
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/user/login", bytes.NewBuffer([]byte(inpBody)))
	handler.Login(c)

	assert.Equal(t, wantStatusCode, c.Writer.Status())
	assert.JSONEq(t, wantBody, w.Body.String())

}

func TestGetUser(t *testing.T) {
	wantBody := `{"data": {"id": "1", "name": "moritz", "email": "moritz.e50@gmail.com", "betaUser": false}, "message": "Success"}`
	wantStatusCode := 200

	usecaseMock := &UserUsecaseMock{}
	conf, _ := config.NewConfig()

	handler := Handler{
		logger:      log.New(),
		userUsecase: usecaseMock,
		validator:   validator.NewValidator(),
		config:      conf,
	}
	usecaseMock.On("GetUserByID", mock.Anything).Return(&entity.UserResponseModel{
		ID:       "1",
		Name:     "moritz",
		Email:    "moritz.e50@gmail.com",
		BetaUser: false,
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/user", nil)

	q := c.Request.URL.Query()
	q.Add("userID", "1")
	c.Request.URL.RawQuery = q.Encode()

	handler.GetUser(c)

	assert.Equal(t, wantStatusCode, c.Writer.Status())
	assert.JSONEq(t, wantBody, w.Body.String())

}
