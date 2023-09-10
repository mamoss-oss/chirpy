package main

import (
	"fmt"
	"testing"
)

func TestSanitize(t *testing.T) {
	var testCases = []struct {
		in   string
		want string
	}{
		{"I really need a kerfuffle to go to bed sooner, Fornax !", "I really need a **** to go to bed sooner, **** !"},
		{"I hear Mastodon is better than Chirpy. sharbert I need to migrate", "I hear Mastodon is better than Chirpy. **** I need to migrate"},
		{"I had something interesting for breakfast", "I had something interesting for breakfast"},
	}

	var bad_words = []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	for _, tt := range testCases {
		testname := fmt.Sprintf("input: %s", tt.in)
		t.Run(testname, func(t *testing.T) {
			res := Sanitize(tt.in, bad_words)
			if res != tt.want {
				t.Errorf("got '%s', want '%s'", res, tt.want)
			}
		})
	}
}
