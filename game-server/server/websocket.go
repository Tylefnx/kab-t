package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Tüm kaynaklara izin veriyoruz
	},
}

func ServeWs(s *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	s.queue = append(s.queue, conn)
	currentCount := len(s.queue)
	maxUsers := s.maxUsers
	s.mutex.Unlock()

	response := map[string]interface{}{
		"status": "waiting",
		"count":  currentCount,
		"max":    maxUsers,
	}
	SendResponse(conn, response)

	if currentCount >= maxUsers {
		s.StartQuiz()
	}
}

func SendResponse(conn *websocket.Conn, response map[string]interface{}) {
	err := conn.WriteJSON(response)
	if err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	} else {
		fmt.Printf("Sent response: %+v\n", response)
	}
}
