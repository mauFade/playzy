package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mauFade/playzy/internal/model"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024 // 512KB
)

// Client representa uma conexão cliente
type Client struct {
	manager *Manager
	conn    *websocket.Conn
	send    chan model.Message
	mutex   sync.Mutex
	userID  string

	// New fields for connection management
	lastPing time.Time
	isAlive  bool
}

// ReadPump lê mensagens da conexão WebSocket
func (c *Client) ReadPump() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in ReadPump: %v", r)
		}
		c.manager.unregister <- c
		c.conn.Close()
	}()

	// Configure connection
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.lastPing = time.Now()
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Validate message size
		if len(data) > maxMessageSize {
			log.Printf("Message too large: %d bytes", len(data))
			continue
		}

		// Decodificar JSON
		var message model.Message
		if err := json.Unmarshal(data, &message); err != nil {
			log.Printf("Erro ao decodificar mensagem: %v", err)
			continue
		}

		// Validate message
		if err := c.validateMessage(message); err != nil {
			log.Printf("Invalid message: %v", err)
			continue
		}

		// Garantir que o remetente seja correto
		message.SenderID = c.userID
		message.Timestamp = time.Now()
		message.IsRead = false

		// Tente salvar a mensagem antes de broadcast
		err = c.manager.repository.Create(message)
		if err != nil {
			log.Printf("Erro ao salvar mensagem no banco: %v", err)
			// Continue mesmo com erro para não interromper o fluxo
		}

		// Enviar a mensagem para processamento
		select {
		case c.manager.broadcast <- message:
			// Message sent successfully
		default:
			log.Printf("Broadcast channel full, message dropped")
		}
	}
}

// WritePump envia mensagens para a conexão WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in WritePump: %v", r)
		}
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The manager closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

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

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// validateMessage validates the message before processing
func (c *Client) validateMessage(msg model.Message) error {
	if msg.Content == "" {
		return fmt.Errorf("empty message content")
	}
	if msg.ReceiverID == "" {
		return fmt.Errorf("missing receiver ID")
	}
	return nil
}
