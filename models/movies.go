package models

type Film struct {
	// iD       int
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Director string  `json:"director"`
	Rating   float64 `json:"rating"`
}

func NewMovie(name, director string, rating float64) *Film {
	return &Film{
		Name:     name,
		Director: director,
		Rating:   rating,
	}
}

type films = []Film
