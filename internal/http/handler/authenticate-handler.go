package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/mauFade/playzy/internal/repository"
	"github.com/mauFade/playzy/internal/usecase/user"
)

type authenticatePayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserHandler struct {
	db *sql.DB
}

func NewAuthenticateUserHandler(d *sql.DB) *AuthenticateUserHandler {
	return &AuthenticateUserHandler{
		db: d,
	}
}

func (h *AuthenticateUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req authenticatePayload
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	w.Header().Set("Content-Type", "application/json")

	decoder.Decode(&req)

	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "missing required fields"})

		return
	}

	usecase := user.NewAuthenticateUserUseCase(repository.NewUserRepository(h.db))

	res, err := usecase.Execute(&user.AuthenticateRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
