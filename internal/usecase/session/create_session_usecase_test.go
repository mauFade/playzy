package session_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/mauFade/playzy/internal/usecase/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockCreateSessionRepository struct {
	mock.Mock
}

func (m *MockCreateSessionRepository) Create(s *model.SessionModel) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockCreateSessionRepository) FindByID(id uuid.UUID) (*model.SessionModel, error) {
	args := m.Called(id)

	return args.Get(0).(*model.SessionModel), args.Error(1)
}

func (m *MockCreateSessionRepository) FindAvailable(page int) ([]model.SessionModel, error) {
	args := m.Called(page)

	return args.Get(0).([]model.SessionModel), args.Error(1)
}

type MockSessionUserRepository struct {
	mock.Mock
}

func (m *MockSessionUserRepository) FindByEmail(email string) (*model.UserModel, error) {
	args := m.Called(email)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockSessionUserRepository) FindByPhone(phone string) (*model.UserModel, error) {
	args := m.Called(phone)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockSessionUserRepository) FindByGamertag(gamertag string) (*model.UserModel, error) {
	args := m.Called(gamertag)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockSessionUserRepository) FindByID(id string) (*model.UserModel, error) {
	args := m.Called(id)

	return args.Get(0).(*model.UserModel), args.Error(1)
}

func (m *MockSessionUserRepository) Create(user *model.UserModel) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestCreateSessionUseCaseExecuteSuccess(t *testing.T) {
	sr := new(MockCreateSessionRepository)
	ur := new(MockSessionUserRepository)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 6)

	userID := uuid.New()

	existingUser := &model.UserModel{
		ID:        userID,
		Name:      "Existing User",
		Email:     "test@example.com",
		Phone:     "1234567890",
		Gamertag:  "gamer123",
		Password:  string(hash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ur.On("FindByID", userID.String()).Return(existingUser, nil).Once()
	sr.On("Create", mock.Anything).Return(nil).Once()

	uc := session.NewCreateSessionUseCase(sr, ur)

	rank := "Dima"

	res, err := uc.Execute(&session.CreateSessionRequest{
		UserID:    userID.String(),
		Game:      "Game",
		Objective: "Obj",
		Rank:      &rank,
		IsRanked:  true,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.UserID.String(), userID.String())
}
