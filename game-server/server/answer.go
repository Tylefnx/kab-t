package server

import "fmt"

func (s *Server) calculateScores(question map[string]interface{}) {
	correctAnswer := question["answer"].(int)
	if answers, exists := s.answers[question["id"].(int)]; exists {
		for playerID, answer := range answers {
			fmt.Printf("Player %s answer: %d, correct answer: %d\n", playerID, answer, correctAnswer) // Hata ayıklama için loglayalım
			if answer == correctAnswer {
				if _, exists := s.scores[playerID]; !exists {
					s.scores[playerID] = 0
				}
				s.scores[playerID] += 1
			}
			fmt.Printf("Score for player %s: %d\n", playerID, s.scores[playerID]) // Puanları loglayalım
		}
	} else {
		fmt.Printf("No answers found for question %d\n", question["id"].(int)) // Hata ayıklama için loglayalım
	}
	fmt.Printf("Scores after question %d: %+v\n", question["id"].(int), s.scores) // Puanları loglayalım
}
