package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mauFade/playzy/internal/constants"
	"github.com/mauFade/playzy/internal/model"
)

func (m *Manager) ServeWs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID é obrigatório", http.StatusBadRequest)
		return
	}

	m.mutex.Lock()
	if existingClient, ok := m.clients[userID]; ok {
		// Desconectar cliente existente para permitir a nova conexão
		existingClient.conn.Close()
		delete(m.clients, userID)
	}
	m.mutex.Unlock()

	// Fazer upgrade da conexão HTTP para WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erro ao fazer upgrade para WebSocket: %v", err)
		return
	}

	// Criar um novo cliente
	client := &Client{
		manager: m,
		conn:    conn,
		send:    make(chan model.Message, 256),
		userID:  userID,
	}

	// Registrar o cliente
	m.register <- client

	// Iniciar goroutines para leitura e escrita
	go client.WritePump()
	go client.ReadPump()
}

func (m *Manager) GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	// Obter ID do usuário do contexto
	userID, ok := r.Context().Value(constants.UserKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}

	// Obter ID do outro usuário
	otherUserID := r.URL.Query().Get("otherUserId")
	if otherUserID == "" {
		http.Error(w, "otherUserId é obrigatório", http.StatusBadRequest)
		return
	}

	// Obter limite de mensagens
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Obter offset (paginação)
	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Buscar mensagens do banco de dados
	query := `
			SELECT id, content, user_id, receiver_id, created_at, is_read
			FROM messages
			WHERE (user_id = $1 AND receiver_id = $2) OR (user_id = $2 AND receiver_id = $1)
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
	`

	rows, err := m.db.Query(query, userID, otherUserID, limit, offset)
	if err != nil {
		http.Error(w, "Erro ao buscar mensagens", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type MessageResponse struct {
		ID         string    `json:"id"`
		Content    string    `json:"content"`
		SenderID   string    `json:"senderId"`
		ReceiverID string    `json:"receiverId"`
		Timestamp  time.Time `json:"timestamp"`
		IsRead     bool      `json:"isRead"`
		IsMine     bool      `json:"isMine"`
	}

	var messages []MessageResponse
	for rows.Next() {
		var msg MessageResponse
		err := rows.Scan(&msg.ID, &msg.Content, &msg.SenderID, &msg.ReceiverID, &msg.Timestamp, &msg.IsRead)
		if err != nil {
			continue
		}

		// Marcar se a mensagem é do usuário atual
		msg.IsMine = msg.SenderID == userID

		messages = append(messages, msg)
	}

	// Inverter a ordem para mostrar as mais antigas primeiro
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// Marcar mensagens como lidas (se o usuário for o destinatário)
	if len(messages) > 0 {
		_, err := m.db.Exec(`
					UPDATE messages
					SET is_read = true
					WHERE receiver_id = $1 AND user_id = $2 AND is_read = false
			`, userID, otherUserID)

		if err != nil {
			log.Printf("Erro ao marcar mensagens como lidas: %v", err)
		}
	}

	// Enviar como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
