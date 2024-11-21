package user_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/mauFade/playzy/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthUserRepository struct {
	mock.Mock
}

func (m *MockAuthUserRepository) FindByEmail(email string) (*model.UserModel, error) {
	args := m.Called(email)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockAuthUserRepository) FindByPhone(phone string) (*model.UserModel, error) {
	args := m.Called(phone)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockAuthUserRepository) FindByGamertag(gamertag string) (*model.UserModel, error) {
	args := m.Called(gamertag)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockAuthUserRepository) Create(user *model.UserModel) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestAuthenticateUserUseCaseExecuteSuccess(t *testing.T) {
	mockRepo := new(MockAuthUserRepository)
	useCase := user.NewAuthenticateUserService(mockRepo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 6)

	existingUser := &model.UserModel{
		ID:        uuid.New(),
		Name:      "Existing User",
		Email:     "test@example.com",
		Phone:     "1234567890",
		Gamertag:  "gamer123",
		Password:  string(hash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil).Once()

	request := &user.AuthenticateRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	res, err := useCase.Execute(request)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.UserID)
	assert.NotNil(t, res.Token)

	mockRepo.AssertExpectations(t)
}
