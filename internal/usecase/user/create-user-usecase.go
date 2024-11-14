package user

import "github.com/mauFade/playzy/internal/repository"

type CreateUserUseCase struct {
	userRepository *repository.UserRepository
}

func NewCreateUserUseCase(r *repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: r,
	}
}

func (uc *CreateUserUseCase) Execute() {}
