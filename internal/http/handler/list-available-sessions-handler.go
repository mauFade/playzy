package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mauFade/playzy/internal/repository"
	"github.com/mauFade/playzy/internal/usecase/session"
)

type ListAvailableSessionsHandler struct {
	db *sql.DB
}

func NewListAvailableSessionsHandler(d *sql.DB) *ListAvailableSessionsHandler {
	return &ListAvailableSessionsHandler{
		db: d,
	}
}

func (h *ListAvailableSessionsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page := r.URL.Query().Get("page")
	pNum, err := strconv.Atoi(page)

	fmt.Println(page, pNum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})

		return
	}

	sr := repository.NewSessionRepository(h.db)

	uc := session.NewListAvailableSessionsUseCase(sr)

	resp, err := uc.Execute(&session.ListAvailableSessionsRequest{
		Page: pNum,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
