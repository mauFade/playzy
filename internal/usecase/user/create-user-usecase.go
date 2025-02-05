package user

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type CreateUserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Gamertag  string     `json:"gamertag"`
	Token     string     `json:"token"`
	Deleted   bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func NewCreateUserUseCase(r repository.UserRepositoryInterface) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: r,
	}
}

func (uc *CreateUserUseCase) Execute(data *CreateUserRequest) (*CreateUserResponse, error) {
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

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.GetID().String()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, err
	}

	err = uc.userRepository.Create(user)

	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		ID:        user.GetID(),
		Name:      user.GetName(),
		Email:     user.GetEmail(),
		Phone:     user.GetPhone(),
		Gamertag:  user.GetGamertag(),
		Token:     tokenString,
		Deleted:   user.IsDeleted(),
		DeletedAt: user.GetDeletedAt(),
		UpdatedAt: user.GetUpdatedAt(),
		CreatedAt: user.GetCreatedAt(),
	}, nil
}
