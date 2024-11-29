package model

import (
	"time"

	"github.com/google/uuid"
)

type SessionModel struct {
	ID        uuid.UUID `json:"id"`         // type:uuid
	Game      string    `json:"game"`       // type:varchar
	UserID    uuid.UUID `json:"user_id"`    // type:uuid
	Objective string    `json:"objetive"`   // type:varchar
	Rank      *string   `json:"rank"`       // type:varchar nullable:true
	IsRanked  bool      `json:"is_ranked"`  // type:bool
	UpdatedAt time.Time `json:"updated_at"` // type:timestamp
	CreatedAt time.Time `json:"created_at"` // type:timestamp
}

func NewSessionModel(
	id, userId uuid.UUID, game, obj string, rank *string, isRanked bool, updatedAt, createdAt time.Time,
) *SessionModel {
	return &SessionModel{
		ID:        id,
		Game:      game,
		UserID:    userId,
		Objective: obj,
		Rank:      rank,
		IsRanked:  isRanked,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}
}

func (s *SessionModel) GetID() uuid.UUID {
	return s.ID
}

func (s *SessionModel) GetUserID() uuid.UUID {
	return s.UserID
}

func (s *SessionModel) SetUserID(id uuid.UUID) {
	s.UserID = id
}

func (s *SessionModel) GetGame() string {
	return s.Game
}

func (s *SessionModel) SetGame(g string) {
	s.Game = g
}

func (s *SessionModel) GetObjective() string {
	return s.Objective
}

func (s *SessionModel) SetObjective(o string) {
	s.Objective = o
}

func (s *SessionModel) GetRank() *string {
	return s.Rank
}

func (s *SessionModel) SetRank(r *string) {
	s.Rank = r
}

func (s *SessionModel) GetIsRanked() bool {
	return s.IsRanked
}

func (s *SessionModel) SetIsRankedOrNot() {
	s.IsRanked = !s.IsRanked
}

func (s *SessionModel) GetUpdatedAt() time.Time {
	return s.UpdatedAt
}

func (s *SessionModel) GetCreatedAt() time.Time {
	return s.CreatedAt
}
