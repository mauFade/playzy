package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mauFade/playzy/internal/repository"
)

type AuthenticateUserService struct {
	userRepository repository.UserRepositoryInterface
}

type AuthenticateRequest struct {
	Email    string
	Password string
}

type authenticateResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func (uc *AuthenticateUserService) Execute(data *AuthenticateRequest) (*authenticateResponse, error) {
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

	tokenString, err := token.SignedString([]byte("JWT_SECRET"))

	if err != nil {
		return nil, err
	}

	return &authenticateResponse{
		UserID: user.GetID().String(),
		Token:  tokenString,
	}, nil
}
