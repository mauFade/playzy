package websocket

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mauFade/playzy/internal/model"
	"github.com/mauFade/playzy/internal/repository"
)

// Manager gerencia todas as conexões WebSocket
type Manager struct {
	clients    map[string]*Client
	broadcast  chan model.Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
	db         *sql.DB
	repository repository.MessageRepositoryInterface

	rateLimiter map[string]time.Time
}

// Configuração do upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Replace with your actual origin checking logic
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:3000", "https://yourdomain.com"}
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true
			}
		}
		return false
	},
	// Add handshake timeout
	HandshakeTimeout: 10 * time.Second,
}

// NewManager cria um novo gerenciador
func NewManager(db *sql.DB, repo repository.MessageRepositoryInterface) *Manager {
	return &Manager{
		clients:     make(map[string]*Client),
		broadcast:   make(chan model.Message, 1000), // Increased buffer size
		register:    make(chan *Client, 100),
		unregister:  make(chan *Client, 100),
		mutex:       sync.RWMutex{},
		db:          db,
		repository:  repo,
		rateLimiter: make(map[string]time.Time),
	}
}

// Start inicia o gerenciador em uma goroutine separada
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			// Registrar um novo cliente pelo userID
			m.mutex.Lock()
			m.clients[client.userID] = client
			m.mutex.Unlock()
			log.Printf("Cliente %s conectado. Total: %d", client.userID, len(m.clients))

		case client := <-m.unregister:
			// Desregistrar um cliente
			m.mutex.Lock()
			if _, ok := m.clients[client.userID]; ok {
				delete(m.clients, client.userID)
				close(client.send)
			}
			m.mutex.Unlock()
			log.Printf("Cliente %s desconectado. Total: %d", client.userID, len(m.clients))

		case message := <-m.broadcast:
			// Enviar mensagem para todos os clientes
			m.mutex.Lock()
			for _, client := range m.clients {
				select {
				case client.send <- message:
					// Mensagem enviada com sucesso
				default:
					// Cliente não consegue receber mensagens
					close(client.send)
					delete(m.clients, client.userID)
				}
			}
			m.mutex.Unlock()
		}
	}
}
