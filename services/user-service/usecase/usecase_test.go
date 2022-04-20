package usecase

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/moshrank/spacey-backend/config"
	"github.com/moshrank/spacey-backend/pkg/auth"
	"github.com/moshrank/spacey-backend/services/user-service/entity"
	"github.com/stretchr/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func (m *UserStoreMock) SaveUser(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *UserStoreMock) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserStoreMock) GetUserByID(id string) (*entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		TestName  string
		inpUser   interface{}
		expOutRes *entity.UserResponseModel
		expOutErr error
	}{
		{
			"Valid User",
			&entity.UserReq{
				Name:     "moritz",
				Email:    "moritz.e50@gmail.com",
				Password: "test_password",
			},
			&entity.UserResponseModel{
				ID:       "1",
				Name:     "moritz",
				Email:    "moritz.e50@gmail.com",
				BetaUser: false,
			},
			nil,
		},
	}

	for _, test := range tests {
		storeMock := &UserStoreMock{}
		storeMock.On("SaveUser", mock.Anything).Return("1", nil)

		cfg, _ := config.NewConfig()
		usecase := NewUserUseCase(log.New(), auth.NewJWT(cfg), storeMock)

		res, err := usecase.CreateUser(test.inpUser)

		assert.Equal(t, test.expOutRes, res)
		assert.Equal(t, test.expOutErr, err)

	}
}
