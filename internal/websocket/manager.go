package websocket

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

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
	mutex      sync.Mutex
	db         *sql.DB
	repository repository.MessageRepositoryInterface
}

// NewManager cria um novo gerenciador
func NewManager(db *sql.DB, repo repository.MessageRepositoryInterface) *Manager {
	return &Manager{
		clients:    make(map[string]*Client),
		broadcast:  make(chan model.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      sync.Mutex{},
		db:         db,
		repository: repo,
	}
}

// Configuração do upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
