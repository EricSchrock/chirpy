package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var chirpAPI string = "/api/validate_chirp"
var chirpLengthLimit int = 140
var profanities []string = []string{"kerfuffle", "sharbert", "fornax"}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	if err := decoder.Decode(&req); err != nil {
		respondWithServerError(w, err)
		return
	}

	if len(req.Body) > chirpLengthLimit {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type response struct {
		Body string `json:"cleaned_body"`
	}

	respondWithJSON(w, http.StatusOK, response{filterProfanity(req.Body)})
}

func filterProfanity(chirp string) string {
	words := strings.Split(chirp, " ")
	for i, word := range words {
		for _, profanity := range profanities {
			if strings.ToLower(word) == profanity {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}

func respondWithServerError(w http.ResponseWriter, err error) {
	log.Print(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	type response struct {
		Error string `json:"error"`
	}

	resp := response{message}
	data, err := json.Marshal(resp)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
