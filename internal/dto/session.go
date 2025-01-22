package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionWithUser struct {
	ID           uuid.UUID `json:"id"`
	Game         string    `json:"game"`
	UserID       uuid.UUID `json:"user_id"`
	Objective    string    `json:"objetive"`
	Rank         *string   `json:"rank"`
	IsRanked     bool      `json:"is_ranked"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
	UserName     string    `json:"user_name"`
	UserGamertag string    `json:"user_gamertag"`
	Email        string    `json:"email"`
}

type SessionsPageResponse struct {
	Page       int               `json:"page"`
	TotalPages int               `json:"total_pages"`
	Sessions   []SessionWithUser `json:"sessions"`
}
