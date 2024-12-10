package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	server *Server
}

type Answer struct {
	PlayerID   string `json:"player_id"`
	QuestionID int    `json:"question_id"`
	Answer     int    `json:"answer"`
}

func NewHandler(s *Server) *Handler {
	return &Handler{server: s}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Tüm kaynaklara izin veriyoruz
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

func (h *Handler) HandleAnswer(w http.ResponseWriter, r *http.Request) {
	var answer Answer
	err := json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.server.mutex.Lock()
	defer h.server.mutex.Unlock()

	if _, exists := h.server.answers[answer.QuestionID]; !exists {
		h.server.answers[answer.QuestionID] = make(map[string]int)
	}
	h.server.answers[answer.QuestionID][answer.PlayerID] = answer.Answer
	w.WriteHeader(http.StatusOK)
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
		fmt.Printf("Sending question to client: %+v\n", response) // Hata ayıklama için
		h.sendResponse(conn, response)
	}
}

func (h *Handler) ShowCorrectAnswer(question map[string]interface{}) {
	h.server.mutex.Lock()
	defer h.server.mutex.Unlock()

	for _, conn := range h.server.queue {
		response := map[string]interface{}{
			"status":         "answer",
			"correct_answer": question["answer"],
		}
		h.sendResponse(conn, response)
	}
}
