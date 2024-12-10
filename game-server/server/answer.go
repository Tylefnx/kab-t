package server

func (s *Server) calculateScores(question map[string]interface{}) {
	correctAnswer := question["answer"].(int)
	for playerID, answer := range s.answers[question["id"].(int)] {
		if answer == correctAnswer {
			s.scores[playerID] += 1
		}
	}
}
