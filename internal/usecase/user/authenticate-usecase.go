package user

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mauFade/playzy/internal/repository"
)

type AuthenticateUserUseCase struct {
	userRepository repository.UserRepositoryInterface
}

type AuthenticateRequest struct {
	Email    string
	Password string
}

type authenticateResponse struct {
	UserID   string `json:"user_id"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Gamertag string `json:"gamertag"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

func NewAuthenticateUserUseCase(r repository.UserRepositoryInterface) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{
		userRepository: r,
	}
}

func (uc *AuthenticateUserUseCase) Execute(data *AuthenticateRequest) (*authenticateResponse, error) {
	user, err := uc.userRepository.FindByEmail(data.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found with this email")
	}

	err = user.ComparePasswords(data.Password)

	if err != nil {
		return nil, errors.New("wrong password")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.GetID().String()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, err
	}

	return &authenticateResponse{
		UserID:   user.GetID().String(),
		Token:    tokenString,
		Name:     user.GetName(),
		Email:    user.GetEmail(),
		Gamertag: user.GetGamertag(),
		Phone:    user.GetPhone(),
		Avatar:   "https://i.pinimg.com/736x/93/13/4b/93134b183c46201053495d49f953be58.jpg",
	}, nil
}
