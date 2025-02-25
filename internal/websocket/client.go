package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mauFade/playzy/internal/model"
)

// Client representa uma conexão cliente
type Client struct {
	manager *Manager
	conn    *websocket.Conn
	send    chan model.Message
	mutex   sync.Mutex
	userID  string
}

// ReadPump lê mensagens da conexão WebSocket
func (c *Client) ReadPump() {
	defer func() {
		c.manager.unregister <- c
		c.conn.Close()
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erro: %v", err)
			}
			break
		}

		// Decodificar JSON
		var message model.Message
		if err := json.Unmarshal(data, &message); err != nil {
			log.Printf("Erro ao decodificar mensagem: %v", err)
			continue
		}

		// Garantir que o remetente seja correto
		message.SenderID = c.userID
		message.Timestamp = time.Now()
		message.IsRead = false

		log.Printf("Mensagem recebida: %+v", message) // Adicione este log

		// Tente salvar a mensagem antes de broadcast
		err = c.manager.repository.Create(message)

		if err != nil {
			log.Printf("Erro ao salvar mensagem no banco: %v", err)
			// Continue mesmo com erro para não interromper o fluxo
		}

		// Enviar a mensagem para processamento
		c.manager.broadcast <- message
	}
}

// WritePump envia mensagens para a conexão WebSocket
func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		// Converter a mensagem para JSON
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Erro ao codificar mensagem: %v", err)
			continue
		}

		c.mutex.Lock()
		err = c.conn.WriteMessage(websocket.TextMessage, data)
		c.mutex.Unlock()

		if err != nil {
			return
		}
	}
}
