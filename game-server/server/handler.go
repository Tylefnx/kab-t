package server

import (
	"encoding/json"
	"fmt"
	"game-server/models"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	server *Server
}

func NewHandler(s *Server) *Handler {
	return &Handler{server: s}
}

func (h *Handler) ServeWs(w http.ResponseWriter, r *http.Request) {
	ServeWs(h.server, w, r)
}
func (h *Handler) HandleAnswer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleAnswer fonksiyonu çağrıldı")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Error reading body: %v\n", err)
		return
	}
	fmt.Printf("REQUEST BODY: %s\n", body)

	var answer models.Answer
	err = json.Unmarshal(body, &answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	fmt.Printf("Decoded answer: %+v\n", answer)

	h.server.mutex.Lock()
	defer h.server.mutex.Unlock()

	if _, exists := h.server.answers[answer.QuestionID]; !exists {
		h.server.answers[answer.QuestionID] = make(map[string]int)
	}
	h.server.answers[answer.QuestionID][answer.PlayerID] = answer.Answer

	fmt.Printf("Received answer: %+v\n", answer)
	fmt.Printf("Current answers map for question %d: %+v\n", answer.QuestionID, h.server.answers[answer.QuestionID])
	fmt.Printf("Total answers map: %+v\n", h.server.answers)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) sendResponse(conn *websocket.Conn, response map[string]interface{}) {
	SendResponse(conn, response)
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
		fmt.Printf("Sending question to client: %+v\n", response)
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
