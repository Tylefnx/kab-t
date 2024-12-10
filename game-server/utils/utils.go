package utils

import "math/rand"

func GenerateRandomNumber() int {
	return rand.Intn(9000) + 1000 // 1000 ile 9999 arasında bir sayı
}
