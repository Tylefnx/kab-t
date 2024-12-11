package server

import (
	"fmt"
	"game-server/models"
	"game-server/question"
	"sort"
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
		time.Sleep(5 * time.Second) // 10 saniye bekleme süresi
		s.handler.ShowCorrectAnswer(question)
		s.calculateScores(question)
		time.Sleep(2 * time.Second) // 2 saniye doğru cevabı gösterme süresi
	}
	s.BroadcastLeaderboard()
	s.queue = []*websocket.Conn{} // Bu satırı en sona taşıyalım
}

func (s *Server) BroadcastLeaderboard() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	scores := []models.Player{}
	for id, score := range s.scores {
		scores = append(scores, models.Player{ID: id, Score: score})
	}

	// Skorları azalan sırada sıralayalım
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	fmt.Printf("Broadcasting leaderboard: %+v\n", scores)

	for _, conn := range s.queue {
		response := map[string]interface{}{
			"status": "leaderboard",
			"scores": scores,
		}
		s.handler.sendResponse(conn, response)
	}
}
