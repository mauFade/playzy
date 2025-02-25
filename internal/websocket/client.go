package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client representa uma conexão cliente
type Client struct {
	manager *Manager
	conn    *websocket.Conn
	send    chan []byte
	mutex   sync.Mutex
}

// ReadPump lê mensagens da conexão WebSocket
func (c *Client) ReadPump() {
	defer func() {
		c.manager.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erro: %v", err)
			}
			break
		}

		// Enviar a mensagem para todos
		c.manager.broadcast <- message
	}
}

// WritePump envia mensagens para a conexão WebSocket
func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		c.mutex.Lock()
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		c.mutex.Unlock()

		if err != nil {
			return
		}
	}
}
