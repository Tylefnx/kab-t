package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	server *Server
}

func NewHandler(s *Server) *Handler {
	return &Handler{server: s}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Tüm kaynaklara izin veriyoruz
		return true
	},
}

func (h *Handler) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	h.server.mutex.Lock()
	h.server.queue = append(h.server.queue, conn)
	currentCount := len(h.server.queue)
	maxUsers := h.server.maxUsers
	h.server.mutex.Unlock()

	response := map[string]interface{}{
		"status": "waiting",
		"count":  currentCount,
		"max":    maxUsers,
	}
	h.sendResponse(conn, response)

	if currentCount >= maxUsers {
		h.server.StartQuiz()
	}
}

func (h *Handler) sendResponse(conn *websocket.Conn, response map[string]interface{}) {
	conn.WriteJSON(response)
}

func (h *Handler) BroadcastQuestion(question map[string]interface{}) {
	h.server.mutex.Lock()
	defer h.server.mutex.Unlock()

	for _, conn := range h.server.queue {
		response := map[string]interface{}{
			"status":   "quiz",
			"question": question,
			"timeout":  10,
		}
		h.sendResponse(conn, response)
	}
}
