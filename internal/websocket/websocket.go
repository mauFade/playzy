package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Configure o upgrader para converter requisições HTTP em conexões WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite qualquer origem em desenvolvimento, ajuste em produção
	},
}

// Client representa uma conexão de cliente WebSocket
type Client struct {
	ID       string
	UserID   string
	RoomID   string
	Username string
	Conn     *websocket.Conn
	Pool     *Pool
	mu       sync.Mutex
}

// Pool gerencia todos os clientes conectados
type Pool struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	mu         sync.RWMutex
}

// Message define a estrutura da mensagem
type Message struct {
	ID        string    `json:"id,omitempty"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	RoomID    string    `json:"roomId"`
	UserID    string    `json:"userId,omitempty"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
}

// NewPool cria um novo pool
func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}
