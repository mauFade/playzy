package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/mauFade/playzy/internal/constants"
	"github.com/mauFade/playzy/internal/repository"
	"github.com/mauFade/playzy/internal/usecase/session"
)

type CreateSessionHandler struct {
	db *sql.DB
}

type createSessionRequest struct {
	Game      string  `json:"game"`
	Objective string  `json:"objective"`
	Rank      *string `json:"rank"`
	IsRanked  bool    `json:"is_ranked"`
}

func NewCreateSessionHandler(d *sql.DB) *CreateSessionHandler {
	return &CreateSessionHandler{
		db: d,
	}
}

func (h *CreateSessionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.UserKey).(string)

	var req createSessionRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	w.Header().Set("Content-Type", "application/json")

	decoder.Decode(&req)

	if req.Game == "" || req.Objective == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "missing required fields"})

		return
	}

	sr := repository.NewSessionRepository(h.db)
	ur := repository.NewUserRepository(h.db)

	usecase := session.NewCreateSessionUseCase(sr, ur)

	response, err := usecase.Execute(&session.CreateSessionRequest{
		UserID:    userID,
		Game:      req.Game,
		Objective: req.Objective,
		Rank:      req.Rank,
		IsRanked:  req.IsRanked,
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
