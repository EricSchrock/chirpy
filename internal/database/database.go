package database

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

var chirps = []Chirp{}

func CreateChirp(body string) (Chirp, error) {
	chirp := Chirp{
		ID:   len(chirps) + 1,
		Body: body,
	}

	chirps = append(chirps, chirp)

	return chirp, nil
}

func GetChirps() ([]Chirp, error) {
	return chirps, nil
}
