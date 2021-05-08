package main

import "math/rand"

func GenerateRandomString(length int, characters string) string {
	letterRunes := []rune(characters)
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}