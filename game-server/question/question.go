package question

import "fmt"

func GenerateQuestions(n int) []map[string]interface{} {
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
	}
	return questions
}
