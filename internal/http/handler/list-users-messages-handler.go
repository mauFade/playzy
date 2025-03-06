package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mauFade/playzy/internal/constants"
	"github.com/mauFade/playzy/internal/repository"
	"github.com/mauFade/playzy/internal/usecase/message"
)

type ListUsersMessagesHandler struct {
	db *sql.DB
}

type MessageResponse struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	SenderID   string    `json:"senderId"`
	ReceiverID string    `json:"receiverId"`
	Timestamp  time.Time `json:"timestamp"`
	IsRead     bool      `json:"isRead"`
	IsMine     bool      `json:"isMine"`
}

func NewListUsersMessagesHandler(d *sql.DB) *ListUsersMessagesHandler {
	return &ListUsersMessagesHandler{
		db: d,
	}
}

func (h *ListUsersMessagesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(constants.UserKey).(string)
	otherUserID := r.URL.Query().Get("otherUserID")

	if otherUserID == "" {
		http.Error(w, "O parâmetro otherUserID é obrigatório", http.StatusBadRequest)
		return
	}

	// Obter limite de mensagens (opcional)
	limit := 50 // valor padrão
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Obter offset para paginação (opcional)
	offset := 0 // valor padrão
	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	mr := repository.NewMessageRepository(h.db)
	uc := message.NewListUsersMessagesUseCase(mr)

	messages, err := uc.Execute(&message.ListUsersMessagesRequest{
		UserID:      userID,
		OtherUserID: otherUserID,
		Limit:       limit,
		Offset:      offset,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})

		return
	}

	// Inverter a ordem para mostrar as mensagens mais antigas primeiro
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
