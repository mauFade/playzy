package websocket

import (
	"log"
	"net/http"
	"time"
)

// Start inicia o pool para gerenciar as conexões
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.mu.Lock()
			pool.Clients[client] = true
			pool.mu.Unlock()

			// Notificar que um novo usuário entrou na sala
			pool.Broadcast <- Message{
				Type:      "system",
				Content:   client.Username + " entrou no chat",
				RoomID:    client.RoomID,
				Timestamp: time.Now(),
			}

			log.Printf("Cliente registrado. Total de conexões: %d", len(pool.Clients))

		case client := <-pool.Unregister:
			pool.mu.Lock()
			if _, ok := pool.Clients[client]; ok {
				delete(pool.Clients, client)
				close := client.Conn.Close

				// Evite fechar uma conexão já fechada
				if client.Conn != nil && close != nil {
					_ = client.Conn.Close()
				}
			}
			pool.mu.Unlock()

			// Notificar que o usuário saiu da sala
			if client.Username != "" {
				pool.Broadcast <- Message{
					Type:      "system",
					Content:   client.Username + " saiu do chat",
					RoomID:    client.RoomID,
					Timestamp: time.Now(),
				}
			}

			log.Printf("Cliente removido. Total de conexões: %d", len(pool.Clients))

		case message := <-pool.Broadcast:
			pool.mu.RLock()
			for client := range pool.Clients {
				// Só envie a mensagem para clientes que estão na mesma sala
				if client.RoomID == message.RoomID {
					client.mu.Lock()
					err := client.Conn.WriteJSON(message)
					client.mu.Unlock()

					if err != nil {
						log.Printf("Erro ao enviar mensagem: %v", err)
						client.Conn.Close()
						pool.Unregister <- client
					}
				}
			}
			pool.mu.RUnlock()
		}
	}
}

// ReadMessages processa as mensagens recebidas do cliente
func (c *Client) ReadMessages() {
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
			// Salvar a mensagem no banco de dados aqui
			// saveMessageToDB(msg)

			// Enviar a mensagem para todos na sala
			c.Pool.Broadcast <- msg
		}
	}
}

// ServeWs gerencia a conexão WebSocket
func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros de query
	userID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")
	roomID := r.URL.Query().Get("roomId")

	// Validar parâmetros
	if userID == "" || username == "" || roomID == "" {
		http.Error(w, "Parâmetros userID, username e roomID são obrigatórios", http.StatusBadRequest)
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
		Pool:     pool,
	}

	// Registrar o cliente no pool
	pool.Register <- client

	// Iniciar a leitura de mensagens em uma goroutine
	go client.ReadMessages()
}
