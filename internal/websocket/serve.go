package websocket

import (
	"log"
	"net/http"
)

func (m *Manager) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		manager: m,
		conn:    conn,
		send:    make(chan []byte, 256),
	}

	client.manager.register <- client

	go client.WritePump()
	go client.ReadPump()
}
