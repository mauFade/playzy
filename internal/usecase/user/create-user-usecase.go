package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/mauFade/playzy/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase struct {
	userRepository repository.UserRepositoryInterface
}

type CreateUserRequest struct {
	Name     string
	Email    string
	Phone    string
	Password string
	Gamertag string
}

func NewCreateUserUseCase(r repository.UserRepositoryInterface) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: r,
	}
}

func (uc *CreateUserUseCase) Execute(data *CreateUserRequest) (*model.UserModel, error) {
	emailExists, err := uc.userRepository.FindByEmail(data.Email)

	if err != nil {
		return nil, err
	}

	if emailExists != nil {
		return nil, errors.New("this email is already in use")
	}

	phoneExists, err := uc.userRepository.FindByPhone(data.Phone)

	if err != nil {
		return nil, err
	}

	if phoneExists != nil {
		return nil, errors.New("this phone is already in use")
	}

	gamertagExists, err := uc.userRepository.FindByGamertag(data.Gamertag)

	if err != nil {
		return nil, err
	}

	if gamertagExists != nil {
		return nil, errors.New("this gamertag is already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 6)

	if err != nil {
		return nil, err
	}

	user := model.NewUserModel(
		uuid.New(),
		data.Name,
		data.Email,
		data.Phone,
		data.Gamertag,
		string(hash),
		false,
		nil,
		time.Now(),
		time.Now(),
	)

	err = uc.userRepository.Create(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
