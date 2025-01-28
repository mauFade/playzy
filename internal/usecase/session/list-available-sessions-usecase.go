package session

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauFade/playzy/internal/repository"
)

type ListAvailableSessionsUseCase struct {
	sr repository.SessionRepositoryInterface
}

type ListAvailableSessionsRequest struct {
	Page int
}

type UserData struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Gamertag string    `json:"gamertag"`
	Avatar   string    `json:"avatar"`
}

type AvailableSessionsResponse struct {
	ID        uuid.UUID `json:"id"`
	Game      string    `json:"game"`
	Objective string    `json:"objetive"`
	Rank      *string   `json:"rank"`
	IsRanked  bool      `json:"is_ranked"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	User      UserData  `json:"user"`
}

type SessionsPageResponse struct {
	Page       int                         `json:"page"`
	TotalPages int                         `json:"total_pages"`
	Sessions   []AvailableSessionsResponse `json:"sessions"`
}

func NewListAvailableSessionsUseCase(s repository.SessionRepositoryInterface) *ListAvailableSessionsUseCase {
	return &ListAvailableSessionsUseCase{
		sr: s,
	}
}

func (u *ListAvailableSessionsUseCase) Execute(data *ListAvailableSessionsRequest) (*SessionsPageResponse, error) {
	sessions, err := u.sr.FindAvailable(data.Page)

	if err != nil {
		return nil, err
	}

	var resSessions []AvailableSessionsResponse

	for _, s := range sessions.Sessions {
		ses := AvailableSessionsResponse{
			ID:        s.ID,
			Game:      s.Game,
			Objective: s.Objective,
			Rank:      s.Rank,
			IsRanked:  s.IsRanked,
			UpdatedAt: s.UpdatedAt,
			CreatedAt: s.CreatedAt,
			User: UserData{
				ID:       s.UserID,
				Name:     s.UserName,
				Email:    s.Email,
				Gamertag: s.UserGamertag,
				Avatar:   "https://i.pinimg.com/736x/6e/27/e4/6e27e43f5e02954d08e0bd3be06f7242.jpg",
			},
		}

		resSessions = append(resSessions, ses)
	}

	return &SessionsPageResponse{
		Page:       sessions.Page,
		TotalPages: sessions.TotalPages,
		Sessions:   resSessions,
	}, nil
}
