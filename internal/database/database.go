package database

import (
	"golang.org/x/exp/maps"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type database struct {
	Chirps map[int]Chirp `json:"chirps"`
}

var chirps = database{Chirps: map[int]Chirp{}}

func CreateChirp(body string) (Chirp, error) {

	chirp := Chirp{
		ID:   len(chirps.Chirps) + 1,
		Body: body,
	}

	chirps.Chirps[chirp.ID] = chirp

	return chirp, nil
}

func GetChirps() ([]Chirp, error) {
	return maps.Values(chirps.Chirps), nil
}
