package util

import (
	"math/rand"
	"time"
)

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateRandomDate() time.Time {
	start := time.Now().AddDate(-30, 0, 0).Unix()
	end := time.Now().Unix()
	sec := rand.Int63n(end-start) + start
	return time.Unix(sec, 0)
}
