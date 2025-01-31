package main

type film struct {
	// iD       int
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Director string  `json:"director"`
	Rating   float64 `json:"rating"`
}

func newMovie(name, director string, rating float64) *film {
	return &film{
		Name:     name,
		Director: director,
		Rating:   rating,
	}
}

type films = []film
