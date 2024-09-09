package main

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const shortURLLength = 7

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateShortURL() string {
	b := make([]byte, shortURLLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
