package handler

import (
	"database/sql"
	"net/http"

	"github.com/mauFade/playzy/internal/constants"
)

type CreateSessionHandler struct {
	db *sql.DB
}

func NewCreateSessionHandler(d *sql.DB) *CreateSessionHandler {
	return &CreateSessionHandler{
		db: d,
	}
}

func (h *CreateSessionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.UserKey).(string)

	w.Write([]byte(userID))
}
