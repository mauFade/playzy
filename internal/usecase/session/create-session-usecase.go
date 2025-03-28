package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/model"
	"github.com/mauFade/playzy/internal/repository"
)

type CreateSessionUseCase struct {
	sr repository.SessionRepositoryInterface
	ur repository.UserRepositoryInterface
}

type CreateSessionRequest struct {
	UserID    string
	Game      string
	Objective string
	Rank      *string
	IsRanked  bool
}

func NewCreateSessionUseCase(r repository.SessionRepositoryInterface, u repository.UserRepositoryInterface) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sr: r,
		ur: u,
	}
}

func (uc *CreateSessionUseCase) Execute(data *CreateSessionRequest) (*model.SessionModel, error) {
	user, err := uc.ur.FindByID(data.UserID)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found with this id")
	}

	var isRanked = true

	if data.Rank == nil {
		isRanked = false
	}

	session := model.NewSessionModel(
		uuid.New(),
		user.GetID(),
		data.Game,
		data.Objective,
		data.Rank,
		isRanked,
		time.Now(),
		time.Now(),
	)

	err = uc.sr.Create(session)

	if err != nil {
		return nil, err
	}

	return session, nil
}
