package server

import (
	"encoding/json"
	"net"
)

type Handler struct {
	server *Server
}

func NewHandler(s *Server) *Handler {
	return &Handler{server: s}
}

func (h *Handler) HandleConnection(conn net.Conn) {
	h.server.mutex.Lock()
	defer h.server.mutex.Unlock()

	h.server.queue = append(h.server.queue, conn)
	currentCount := len(h.server.queue)
	maxUsers := h.server.maxUsers

	// Kullanıcıya bekleme durumunu bildir
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

func (h *Handler) sendResponse(conn net.Conn, response map[string]interface{}) {
	encoder := json.NewEncoder(conn)
	encoder.Encode(response)
}

func (h *Handler) BroadcastQuestion(question map[string]interface{}) {
	h.server.mutex.Lock()
	defer h.server.mutex.Unlock()

	for _, conn := range h.server.queue {
		response := map[string]interface{}{
			"status":   "quiz",
			"question": question,
			"timeout":  10, // 10 saniye
		}
		h.sendResponse(conn, response)
	}
}
