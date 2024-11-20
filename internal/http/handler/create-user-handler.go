package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mauFade/playzy/internal/repository"
	"github.com/mauFade/playzy/internal/usecase/user"
)

type createUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Gamertag string `json:"gamertag"`
}

type CreateUserHandler struct {
	db *sql.DB
}

func NewCreateUserHandler(d *sql.DB) *CreateUserHandler {
	return &CreateUserHandler{
		db: d,
	}
}

func (h *CreateUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req createUserPayload
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Gamertag == "" || req.Phone == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	repository := repository.NewUserRepository(h.db)
	usecase := user.NewCreateUserUseCase(repository)

	res, err := usecase.Execute(&user.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Gamertag: req.Gamertag,
	})

}
