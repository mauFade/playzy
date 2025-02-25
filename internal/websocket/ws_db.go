package websocket

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

// Estrutura para armazenar dependências
type WebsocketHandler struct {
	DB   *sql.DB
	Pool *Pool
}

// NewWebsocketHandler cria um novo handler de websocket
func NewWebsocketHandler(db *sql.DB) *WebsocketHandler {
	pool := NewPool()
	go pool.Start()

	return &WebsocketHandler{
		DB:   db,
		Pool: pool,
	}
}

// SaveMessage salva uma mensagem no banco de dados
func (h *WebsocketHandler) SaveMessage(msg Message) (string, error) {
	var messageID string

	err := h.DB.QueryRow(`
		INSERT INTO messages (room_id, user_id, content, message_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id
	`, msg.RoomID, msg.UserID, msg.Content, msg.Type, msg.Timestamp).Scan(&messageID)

	if err != nil {
		log.Printf("Erro ao salvar mensagem: %v", err)
		return "", err
	}

	return messageID, nil
}

// GetMessages recupera mensagens de uma sala específica
func (h *WebsocketHandler) GetMessages(roomID string, limit, offset int) ([]Message, error) {
	if limit <= 0 {
		limit = 50 // Valor padrão
	}

	rows, err := h.DB.Query(`
		SELECT m.id, m.content, m.message_type, m.user_id, m.created_at, u.username
		FROM messages m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.room_id = $1 AND m.is_deleted = false
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`, roomID, limit, offset)

	if err != nil {
		log.Printf("Erro ao recuperar mensagens: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var username sql.NullString
		var timestamp time.Time

		err := rows.Scan(&msg.ID, &msg.Content, &msg.Type, &msg.UserID, &timestamp, &username)
		if err != nil {
			log.Printf("Erro ao escanear mensagem: %v", err)
			continue
		}

		msg.RoomID = roomID
		msg.Timestamp = timestamp

		if username.Valid {
			msg.Sender = username.String
		} else {
			msg.Sender = "Usuário Desconhecido"
		}

		messages = append(messages, msg)
	}

	// Inverte a ordem para mostrar as mais antigas primeiro
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// JoinRoom adiciona um usuário a uma sala
func (h *WebsocketHandler) JoinRoom(userID, roomID string) error {
	_, err := h.DB.Exec(`
		INSERT INTO room_participants (room_id, user_id, joined_at, last_read_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (room_id, user_id) 
		DO UPDATE SET last_read_at = NOW()
	`, roomID, userID)

	return err
}

// UpdateLastRead atualiza o timestamp de última leitura
func (h *WebsocketHandler) UpdateLastRead(userID, roomID string) error {
	_, err := h.DB.Exec(`
		UPDATE room_participants 
		SET last_read_at = NOW()
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)

	return err
}

// Modificar o método ReadMessages no Client para incluir salvamento de mensagens
func (c *Client) ReadMessagesWithDB(handler *WebsocketHandler) {
	defer func() {
		c.Pool.Unregister <- c
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			break
		}

		// Garantir que a mensagem tenha os dados corretos
		msg.UserID = c.UserID
		msg.Sender = c.Username
		msg.RoomID = c.RoomID
		msg.Timestamp = time.Now()

		// Processar diferentes tipos de mensagens
		switch msg.Type {
		case "message":
			// Salvar a mensagem no banco de dados
			messageID, err := handler.SaveMessage(msg)
			if err != nil {
				log.Printf("Erro ao salvar mensagem: %v", err)
				continue
			}

			// Adicionar o ID da mensagem antes de enviar
			msg.ID = messageID

			// Enviar a mensagem para todos na sala
			c.Pool.Broadcast <- msg

			// Atualizar o timestamp de última leitura para o remetente
			handler.UpdateLastRead(c.UserID, c.RoomID)
		}
	}
}

// ServeWsWithDB é uma versão modificada de ServeWs que inclui o handler
func (h *WebsocketHandler) ServeWsWithDB(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros de query
	userID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")
	roomID := r.URL.Query().Get("roomId")

	// Validar parâmetros
	if userID == "" || username == "" || roomID == "" {
		http.Error(w, "Parâmetros userID, username e roomID são obrigatórios", http.StatusBadRequest)
		return
	}

	// Registrar entrada na sala
	err := h.JoinRoom(userID, roomID)
	if err != nil {
		log.Printf("Erro ao registrar entrada na sala: %v", err)
		http.Error(w, "Erro ao entrar na sala", http.StatusInternalServerError)
		return
	}

	// Fazer upgrade da conexão HTTP para WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erro ao fazer upgrade para WebSocket: %v", err)
		return
	}

	// Criar um novo cliente
	client := &Client{
		ID:       userID + "_" + time.Now().String(),
		UserID:   userID,
		RoomID:   roomID,
		Username: username,
		Conn:     conn,
		Pool:     h.Pool,
	}

	// Registrar o cliente no pool
	h.Pool.Register <- client

	// Iniciar a leitura de mensagens em uma goroutine
	go client.ReadMessagesWithDB(h)
}
