package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/mauFade/playzy/internal/repository"
)

type CreateUserUseCase struct {
	userRepository *repository.UserRepository
}

type CreateUserRequest struct {
	Name     string
	Email    string
	Phone    string
	Password string
	Gamertag string
}

func NewCreateUserUseCase(r *repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: r,
	}
}

func (uc *CreateUserUseCase) Execute(data *CreateUserRequest) (*model.UserModel, error) {
	userExist, err := uc.userRepository.FindByEmail(data.Email)

	if err != nil {
		return nil, err
	}

	if userExist != nil {
		return nil, errors.New("this email is already in use")
	}

	user := model.NewUserModel(
		uuid.New(),
		data.Name,
		data.Email,
		data.Phone,
		data.Gamertag,
		data.Password,
		false,
		nil,
		time.Now(),
		time.Now(),
	)

	uc.userRepository.Create(user)

	return user, nil
}
