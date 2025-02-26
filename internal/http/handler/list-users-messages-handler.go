package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mauFade/playzy/internal/constants"
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

	// Obtém os IDs dos usuários da query
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

	// Consulta SQL para buscar mensagens entre os dois usuários
	query := `
		SELECT id, content, user_id, receiver_id, created_at, is_read
		FROM messages
		WHERE (user_id = $1 AND receiver_id = $2) OR (user_id = $2 AND receiver_id = $1)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := h.db.Query(query, userID, otherUserID, limit, offset)
	if err != nil {
		log.Printf("Erro ao buscar mensagens: %v", err)
		http.Error(w, "Erro ao buscar mensagens", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Processar as mensagens
	var messages []MessageResponse
	for rows.Next() {
		var msg MessageResponse
		var userIDStr, receiverIDStr string

		err := rows.Scan(&msg.ID, &msg.Content, &userIDStr, &receiverIDStr, &msg.Timestamp, &msg.IsRead)
		if err != nil {
			log.Printf("Erro ao escanear mensagem: %v", err)
			continue
		}

		msg.SenderID = userIDStr
		msg.ReceiverID = receiverIDStr

		// Determina se a mensagem é do usuário atual
		msg.IsMine = (userIDStr == userID)

		messages = append(messages, msg)
	}

	// Verificar se houve erro na iteração
	if err = rows.Err(); err != nil {
		log.Printf("Erro ao iterar mensagens: %v", err)
		http.Error(w, "Erro ao processar mensagens", http.StatusInternalServerError)
		return
	}

	// Inverter a ordem para mostrar as mensagens mais antigas primeiro
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// Atualizar mensagens como lidas (se o usuário for o destinatário)
	go func() {
		_, err := h.db.Exec(`
			UPDATE messages
			SET is_read = true
			WHERE receiver_id = $1 AND user_id = $2 AND is_read = false
		`, userID, otherUserID)

		if err != nil {
			log.Printf("Erro ao marcar mensagens como lidas: %v", err)
		}
	}()

	// Retornar as mensagens como JSON
	json.NewEncoder(w).Encode(messages)
}
