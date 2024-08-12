package database

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
	c := make([]Chirp, 0, len(chirps.Chirps))
	for _, v := range chirps.Chirps { // replace with maps.Values(chirps.Chirps) when Go 1.23 is released
		c = append(c, v)
	}

	return c, nil
}
