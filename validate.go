package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
)

// validate_chirp checks if the received message has the length of 140 or below.
func validate_chirp(w http.ResponseWriter, r *http.Request) {
	var bad_words = []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	type message struct {
		Body string `json:"body"`
	}

	type cleaned_body struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	// json receive part
	decoder := json.NewDecoder(r.Body)
	msg := message{}

	err := decoder.Decode(&msg)

	if err != nil {
		log.Printf("Error decoding message body: %s", err)
		w.WriteHeader(500)
		return
	}

	// validate message
	l := len(msg.Body)
	if l > 0 && l <= 140 {
		resp := cleaned_body{Cleaned_body: Sanitize(msg.Body, bad_words)}
		respondWithJSON(w, 200, resp)

	} else if l > 140 {
		msg := "Chirp is too long"
		respondWithError(w, 400, msg)
	} else {
		msg := "Something went wrong"
		respondWithError(w, 400, msg)
	}

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		fmt.Printf("Server side error: %d", code)
	}
	type err_resp struct {
		Err_msg string `json:"error"`
	}
	respondWithJSON(w, code, err_resp{
		Err_msg: msg,
	})
}

func Sanitize(s string, words []string) string {
	var return_text []string
	toSlice := strings.Split(s, " ")
	for _, word := range toSlice {
		if slices.Contains(words, strings.ToLower(word)) {
			return_text = append(return_text, "****")
		} else {
			return_text = append(return_text, word)
		}
	}
	return strings.Join(return_text, " ")
}
