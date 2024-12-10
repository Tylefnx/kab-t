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
	answers  map[int]map[string]int // Cevapları kaydetmek için
	scores   map[string]int         // Kullanıcı puanlarını saklamak için
	mutex    sync.Mutex
	handler  *Handler
}

func NewServer(maxUsers int) *Server {
	srv := &Server{
		maxUsers: maxUsers,
		queue:    []*websocket.Conn{},
		answers:  make(map[int]map[string]int),
		scores:   make(map[string]int),
	}
	srv.handler = NewHandler(srv)
	return srv
}

// Handler fonksiyonunu ekleyelim
func (s *Server) Handler() *Handler {
	return s.handler
}

func (s *Server) ListenAndServe(addr string) error {
	http.HandleFunc("/ws", s.handler.ServeWs)
	http.HandleFunc("/answer", s.handler.HandleAnswer)
	fmt.Println("Server started on", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) StartQuiz() {
	fmt.Println("Starting the quiz with", len(s.queue), "players.")
	questions := generateQuestions(10)
	for _, question := range questions {
		fmt.Printf("Broadcasting question: %+v\n", question) // Hata ayıklama için
		s.handler.BroadcastQuestion(question)
		time.Sleep(10 * time.Second) // 10 saniye bekleme süresi
		s.handler.ShowCorrectAnswer(question)
		s.calculateScores(question)
		time.Sleep(2 * time.Second) // 2 saniye doğru cevabı gösterme süresi
	}
	s.queue = []*websocket.Conn{}
}

func (s *Server) calculateScores(question map[string]interface{}) {
	correctAnswer := question["answer"].(int)
	for playerID, answer := range s.answers[question["id"].(int)] {
		if answer == correctAnswer {
			s.scores[playerID] += 1
		}
	}
}

func generateQuestions(n int) []map[string]interface{} {
	questions := []map[string]interface{}{}
	for i := 0; i < n; i++ {
		questionText := fmt.Sprintf("What is %d + %d?", i, i+1)
		answer := i + i + 1
		choices := []int{answer, answer + 1, answer - 1, answer + 2}

		question := map[string]interface{}{
			"id":      i,
			"text":    questionText,
			"choices": choices,
			"answer":  answer,
		}
		questions = append(questions, question)
		fmt.Printf("Generated question: %+v\n", question) // Hata ayıklama için
	}
	return questions
}
