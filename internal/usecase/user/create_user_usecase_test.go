package user_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/mauFade/playzy/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*model.UserModel, error) {
	args := m.Called(email)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockUserRepository) FindByPhone(phone string) (*model.UserModel, error) {
	args := m.Called(phone)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockUserRepository) FindByGamertag(gamertag string) (*model.UserModel, error) {
	args := m.Called(gamertag)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*model.UserModel, error) {
	args := m.Called(id)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.UserModel) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestCreateUserUseCaseExecuteSuccess(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := user.NewCreateUserUseCase(mockRepo)

	mockRepo.On("FindByEmail", "test@example.com").Return((*model.UserModel)(nil), nil)
	mockRepo.On("FindByPhone", "1234567890").Return((*model.UserModel)(nil), nil)
	mockRepo.On("FindByGamertag", "gamer123").Return((*model.UserModel)(nil), nil)
	mockRepo.On("Create", mock.Anything).Return(nil).Once()

	request := &user.CreateUserRequest{
		Name:     "John Doe",
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: "password123",
		Gamertag: "gamer123",
	}

	user, err := useCase.Execute(request)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.GetName())
	assert.Equal(t, "test@example.com", user.GetEmail())
	assert.Equal(t, "1234567890", user.GetPhone())
	assert.Equal(t, "gamer123", user.GetGamertag())

	mockRepo.AssertExpectations(t)
}

func TestCreateUserUseCaseExecuteEmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := user.NewCreateUserUseCase(mockRepo)

	existingUser := &model.UserModel{
		ID:        uuid.New(),
		Name:      "Existing User",
		Email:     "test@example.com",
		Phone:     "1234567890",
		Gamertag:  "gamer123",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil).Once()

	request := &user.CreateUserRequest{
		Name:     "John Doe",
		Email:    "test@example.com",
		Phone:    "9876543210",
		Password: "password123",
		Gamertag: "gamer999",
	}

	user, err := useCase.Execute(request)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "this email is already in use")

	mockRepo.AssertExpectations(t)
}

func TestCreateUserUseCaseExecutePhoneAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := user.NewCreateUserUseCase(mockRepo)

	existingUser := &model.UserModel{
		ID:        uuid.New(),
		Name:      "Existing User",
		Email:     "test@example.com",
		Phone:     "1234567890",
		Gamertag:  "gamer123",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", "test@example.com").Return((*model.UserModel)(nil), nil)
	mockRepo.On("FindByPhone", "1234567890").Return(existingUser, nil).Once()

	request := &user.CreateUserRequest{
		Name:     "John Doe",
		Email:    "test@example.com",
		Phone:    "1234567890",
		Password: "password123",
		Gamertag: "gamer999",
	}

	user, err := useCase.Execute(request)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "this phone is already in use")

	mockRepo.AssertExpectations(t)
}

func TestCreateUserUseCaseExecuteGamertagAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	useCase := user.NewCreateUserUseCase(mockRepo)

	existingUser := &model.UserModel{
		ID:        uuid.New(),
		Name:      "Existing User",
		Email:     "test@example.com",
		Phone:     "1234567890",
		Gamertag:  "gamer123",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", "test@example.com").Return((*model.UserModel)(nil), nil)
	mockRepo.On("FindByPhone", "9876543210").Return((*model.UserModel)(nil), nil)
	mockRepo.On("FindByGamertag", "gamer123").Return(existingUser, nil).Once()

	request := &user.CreateUserRequest{
		Name:     "John Doe",
		Email:    "test@example.com",
		Phone:    "9876543210",
		Password: "password123",
		Gamertag: "gamer123",
	}

	user, err := useCase.Execute(request)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "this gamertag is already in use")

	mockRepo.AssertExpectations(t)
}
