package server

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	maxUsers int
	queue    []net.Conn
	mutex    sync.Mutex
	handler  *Handler
}

func NewServer(maxUsers int) *Server {
	srv := &Server{
		maxUsers: maxUsers,
		queue:    []net.Conn{},
	}
	srv.handler = NewHandler(srv)
	return srv
}

func (s *Server) ListenAndServe(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handler.HandleConnection(conn)
	}
}

func (s *Server) StartQuiz() {
	fmt.Println("Starting the quiz with", len(s.queue), "players.")
	questions := generateQuestions(10) // Örnek olarak 10 soru üret
	for _, question := range questions {
		s.handler.BroadcastQuestion(map[string]interface{}{
			"question": question,
		})
		time.Sleep(10 * time.Second) // 10 saniye bekleme süresi
	}
	s.queue = []net.Conn{} // Queue'u sıfırlıyoruz
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
