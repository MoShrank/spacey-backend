package handler

import (
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/moshrank/spacey-backend/pkg/testingutil"
	"github.com/moshrank/spacey-backend/pkg/validator"
	"github.com/moshrank/spacey-backend/services/learning-service/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SessionUsecaseMock struct {
	mock.Mock
}

func (u *SessionUsecaseMock) CreateLearningSession(
	userID string,
	session *entity.LearningSessionCreateReq,
) (string, error) {
	args := u.Called(userID, session)
	return args.String(0), args.Error(1)
}

func (u *SessionUsecaseMock) FinishLearningSession(
	userID string,
	session *entity.LearningSessionUpdateReq,
) error {
	args := u.Called(userID, session)
	return args.Error(0)
}

func TestCreateInvalidLearningSession(t *testing.T) {
	tests := []struct {
		testName       string
		body           string
		userID         string
		wantStatusCode int
	}{
		{
			"Empty Body",
			``,
			"test_user_id",
			400,
		},
		{
			"Missing DeckID",
			`{"startedAt": "2012-04-23T18:25:43.511Z"}`,
			"test_user_id",
			400,
		},
		{
			"Missing StartedAt",
			`{"deckID": "1"}`,
			"test_user_id",
			400,
		},
		{
			"Missing UserID",
			`{"deckID": "1", "startedAt": "2012-04-23T18:25:43.511Z"}`,
			"",
			400,
		},
	}

	sessionUsecaseMock := new(SessionUsecaseMock)
	validatorObj := validator.NewValidator()

	handler := NewLearningSessionHandler(sessionUsecaseMock, log.New(), validatorObj)

	sessionUsecaseMock.On("CreateLearningSession", mock.Anything, mock.Anything).Return("1", nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c := testingutil.NewTestingContext(w, "POST", "/session", test.body).
				AddQueryParameter("userID", test.userID)

			handler.CreateLearningSession(c.Context)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestCreateValidLearningSession(t *testing.T) {

	body := `{"deckID": "1", "startedAt": "2012-04-23T18:25:43.511Z"}`
	wantStatusCode := 201
	wantBody := `{"message": "Created", "data": {"id": "1"}}`

	sessionUsecaseMock := new(SessionUsecaseMock)
	validatorObj := validator.NewValidator()

	handler := NewLearningSessionHandler(sessionUsecaseMock, log.New(), validatorObj)

	sessionUsecaseMock.On("CreateLearningSession", mock.Anything, mock.Anything).Return("1", nil)

	w := httptest.NewRecorder()
	c := testingutil.NewTestingContext(w, "POST", "/session", body).
		AddQueryParameter("userID", "1")

	handler.CreateLearningSession(c.Context)

	assert.Equal(t, wantStatusCode, c.Writer.Status())
	assert.JSONEq(t, wantBody, w.Body.String())
}

func TestFinishInvalidLearningSession(t *testing.T) {
	tests := []struct {
		testName       string
		body           string
		userID         string
		wantStatusCode int
	}{
		{
			"Empty Body",
			``,
			"test_user_id",
			400,
		},
		{
			"Missing UserID",
			`{"deckID": "1", "startedAt": "2012-04-23T18:25:43.511Z"}`,
			"",
			400,
		},
		{
			"Missing ID",
			`{"finishedAt": ""2012-04-23T18:25:43.511Z"}`,
			"test_user_id",
			400,
		},
		{
			"Missing FinishedAt",
			`{"ID": "1"}`,
			"test_user_id",
			400,
		},
	}

	sessionUsecaseMock := new(SessionUsecaseMock)
	validatorObj := validator.NewValidator()

	handler := NewLearningSessionHandler(sessionUsecaseMock, log.New(), validatorObj)

	sessionUsecaseMock.On("FinishLearningSession", mock.Anything, mock.Anything).Return(nil)

	for _, test := range tests {

		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			c := testingutil.NewTestingContext(w, "PUT", "/session", test.body).
				AddQueryParameter("userID", test.userID)

			handler.FinishLearningSession(c.Context)

			assert.Equal(t, test.wantStatusCode, c.Writer.Status())
		})
	}
}

func TestFinishValidLearningSession(t *testing.T) {
	body := `{"ID": "1", "finishedAt": "2012-04-23T18:25:43.511Z"}`
	wantStatusCode := 200
	wantBody := `{"message": "Success"}`

	sessionUsecaseMock := new(SessionUsecaseMock)
	validatorObj := validator.NewValidator()

	handler := NewLearningSessionHandler(sessionUsecaseMock, log.New(), validatorObj)

	sessionUsecaseMock.On("FinishLearningSession", mock.Anything, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	c := testingutil.NewTestingContext(w, "PUT", "/session", body).
		AddQueryParameter("userID", "1")

	handler.FinishLearningSession(c.Context)

	assert.Equal(t, wantStatusCode, c.Writer.Status())
	assert.JSONEq(t, wantBody, w.Body.String())
}
