package server

import (
	"fmt"
	"game-server/question"
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

func (s *Server) Handler() *Handler {
	return s.handler
}

func (s *Server) StartQuiz() {
	fmt.Println("Starting the quiz with", len(s.queue), "players.")
	questions := question.GenerateQuestions(10)
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
