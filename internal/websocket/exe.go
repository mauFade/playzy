package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Configuração do upgrader para converter requisições HTTP em conexões WebSocket
var u = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite qualquer origem em desenvolvimento
	},
}

// Handler básico para WebSocket
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade da conexão HTTP para WebSocket
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erro ao fazer upgrade para WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Loop simples para receber e enviar mensagens de eco
	for {
		// Lê uma mensagem do cliente
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			break
		}

		// Loga a mensagem recebida
		log.Printf("Mensagem recebida: %s", p)

		// Envia a mesma mensagem de volta (eco)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Printf("Erro ao enviar mensagem: %v", err)
			break
		}
	}
}
