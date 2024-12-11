package utils

import "math/rand"

func GenerateRandomNumber() int {
	return rand.Intn(9000) + 1000
}
