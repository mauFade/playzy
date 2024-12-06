package session

import (
	"github.com/mauFade/playzy/internal/model"
	"github.com/mauFade/playzy/internal/repository"
)

type ListAvailableSessionsUseCase struct {
	sr repository.SessionRepositoryInterface
}

type ListAvailableSessionsRequest struct {
	Page int
}

func NewListAvailableSessionsUseCase(s repository.SessionRepositoryInterface) *ListAvailableSessionsUseCase {
	return &ListAvailableSessionsUseCase{
		sr: s,
	}
}

func (u *ListAvailableSessionsUseCase) Execute(data *ListAvailableSessionsRequest) ([]model.SessionModel, error) {
	sessions, err := u.sr.FindAvailable(data.Page)

	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return []model.SessionModel{}, nil
	}

	return sessions, nil
}
