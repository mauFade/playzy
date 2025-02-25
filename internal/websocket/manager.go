package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Manager gerencia todas as conexões WebSocket
type Manager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

// NewManager cria um novo gerenciador
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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

func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			// Registrar um novo cliente
			m.mutex.Lock()
			m.clients[client] = true
			m.mutex.Unlock()
			log.Printf("Cliente conectado. Total: %d", len(m.clients))

		case client := <-m.unregister:
			// Desregistrar um cliente
			m.mutex.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}
			m.mutex.Unlock()
			log.Printf("Cliente desconectado. Total: %d", len(m.clients))

		case message := <-m.broadcast:
			// Enviar mensagem para todos os clientes
			m.mutex.Lock()
			for client := range m.clients {
				select {
				case client.send <- message:
					// Mensagem enviada com sucesso
				default:
					// Cliente não consegue receber mensagens
					close(client.send)
					delete(m.clients, client)
				}
			}
			m.mutex.Unlock()
		}
	}
}
