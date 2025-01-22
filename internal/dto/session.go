package dto

import "github.com/mauFade/playzy/internal/model"

type SessionsPageResponse struct {
	Page       int                  `json:"page"`
	TotalPages int                  `json:"total_pages"`
	Sessions   []model.SessionModel `json:"sessions"`
}
