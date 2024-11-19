package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type createUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Gamertag string `json:"gamertag"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req createUserPayload
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Fazer validações adicionais, se necessário
	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

}
