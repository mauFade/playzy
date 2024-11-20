package handler

import (
	"database/sql"
	"encoding/json"
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

	w.Header().Set("Content-Type", "application/json")

	decoder.Decode(&req)

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Gamertag == "" || req.Phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "missing required fields fdp"})

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

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

}
