package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	maxUsers int
	queue    []*websocket.Conn
	mutex    sync.Mutex
	handler  *Handler
}

func NewServer(maxUsers int) *Server {
	srv := &Server{
		maxUsers: maxUsers,
		queue:    []*websocket.Conn{},
	}
	srv.handler = NewHandler(srv)
	return srv
}

func (s *Server) ListenAndServe(addr string) error {
	http.HandleFunc("/ws", s.handler.ServeWs)
	fmt.Println("Server started on", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) StartQuiz() {
	fmt.Println("Starting the quiz with", len(s.queue), "players.")
	questions := generateQuestions(10)
	for _, question := range questions {
		s.handler.BroadcastQuestion(map[string]interface{}{
			"question": question,
		})
		time.Sleep(10 * time.Second)
	}
	s.queue = []*websocket.Conn{}
}

func generateQuestions(n int) []map[string]interface{} {
	questions := []map[string]interface{}{}
	for i := 0; i < n; i++ {
		questions = append(questions, map[string]interface{}{
			"id":     i,
			"text":   fmt.Sprintf("What is %d + %d?", i, i+1),
			"answer": i + i + 1,
		})
	}
	return questions
}
