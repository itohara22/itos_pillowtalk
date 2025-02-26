package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sushi/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HandlerStruct struct {
	DbPool *pgxpool.Pool
	T      *template.Template
}

// func renderTemplates(filename string, res http.ResponseWriter) *template.Template {
// 	tmp, err := template.ParseFiles("templates/" + filename)
// 	if err != nil {
// 		http.Error(res, "can't parse html template", http.StatusInternalServerError)
// 	}
// 	return tmp
// }
//
// func renderTemplatesParseGlob(res http.ResponseWriter) *template.Template {
// 	tmp, err := template.ParseGlob("templates/*.html")
// 	if err != nil {
// 		http.Error(res, "can't parse html template", http.StatusInternalServerError)
// 	}
// 	return tmp
// }

func (h *HandlerStruct) Home(res http.ResponseWriter, req *http.Request) {
	rows, err := h.DbPool.Query(context.Background(), "SELECT * FROM movies;")
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(res, "No movies", http.StatusNotFound)
		}
		http.Error(res, "Something went wrong", http.StatusInternalServerError)
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var movies []models.Film
	for rows.Next() {
		var film models.Film
		err := rows.Scan(&film.ID, &film.Name, &film.Director, &film.Rating)
		// nil deference means accessing a pointer of a variable which is not inintialized yet
		if err != nil {
			http.Error(res, "oh no!", http.StatusInternalServerError)
			fmt.Println(err.Error())
			break
		}
		movies = append(movies, film)
	}

	dataToPasedtoTemplate := map[string]interface{}{"movies": movies, "title": "Home"} // inerface{} is like using any here

	// tpl := renderTemplatesParseGlob(res)
	h.T.ExecuteTemplate(res, "home.html", dataToPasedtoTemplate)
	// tpl.Execute(res, nil) // it will execute the template according to the route and filename
}

func (hand *HandlerStruct) GetMovieWithId(res http.ResponseWriter, req *http.Request) {
	movieId := req.PathValue("id")
	result := hand.DbPool.QueryRow(context.Background(), "SELECT id, name, director, rating FROM movies WHERE id = $1;", movieId)

	var data models.Film
	err := result.Scan(&data.ID, &data.Name, &data.Director, &data.Rating)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	// res.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(res).Encode(data)
	// fmt.Println(req.Header.Get("Accept"))

	// tpl := renderTemplatesParseGlob(res)
	hand.T.ExecuteTemplate(res, "movie.html", data)
}

func (h *HandlerStruct) GetNewMovie(res http.ResponseWriter, req *http.Request) {
	err := h.T.ExecuteTemplate(res, "addMovie.html", map[string]any{"title": "Add new movie"})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HandlerStruct) PostMovie(res http.ResponseWriter, req *http.Request) {
	rating, err := strconv.ParseFloat(req.FormValue("rating"), 64)
	if err != nil {
		http.Error(res, "rating is not valid", http.StatusBadRequest)
		return
	}

	movie := models.NewMovie(req.FormValue("movie"), req.FormValue("director"), rating)
	quryString := "INSERT INTO movies (name, director, rating) VALUES ($1, $2, $3)"
	_, err = h.DbPool.Exec(context.Background(), quryString, movie.Name, movie.Director, movie.Rating)
	// fmt.Println(a) // it shows number of items inserted
	if err != nil {
		http.Error(res, "could not add to db", http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, "/", http.StatusSeeOther) // use proper redirect status codes to make it work
}
