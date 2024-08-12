package database

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

var Chirps = []Chirp{}
