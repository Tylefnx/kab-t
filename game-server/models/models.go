package models

type Answer struct {
	PlayerID   string `json:"player_id"`
	QuestionID int    `json:"question_id"`
	Answer     int    `json:"answer"`
}

type Question struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Choices []int  `json:"choices"`
	Answer  int    `json:"answer"`
}

type Player struct {
	ID    string `json:"id"`
	Score int    `json:"score"`
}
