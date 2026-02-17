package main

import "strings"

var profanity = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

func removeProfanity(msg string) string {
	msgWords := strings.Split(msg, " ")
	cleanedWords := make([]string, len(msgWords))

	for i, word := range msgWords {
		if _, ok := profanity[strings.ToLower(word)]; ok {
			word = "****"
		}

		cleanedWords[i] = word
	}

	cleanedMsg := strings.Join(cleanedWords, " ")

	return cleanedMsg
}
