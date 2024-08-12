package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/EricSchrock/chirpy/internal/database"
)

var ChirpAPI string = "/api/chirps"
var ChirpLengthLimit int = 140
var Profanities = []string{"kerfuffle", "sharbert", "fornax"}

func PostChirpHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	if err := decoder.Decode(&req); err != nil {
		respondWithServerError(w, err)
		return
	}

	if len(req.Body) > ChirpLengthLimit {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := database.CreateChirp(filterProfanity((req.Body)))
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func GetChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := database.GetChirps()
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func filterProfanity(chirp string) string {
	words := strings.Split(chirp, " ")
	for i, word := range words {
		for _, profanity := range Profanities {
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
