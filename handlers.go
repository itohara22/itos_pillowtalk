package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func renderTemplates(filename string, res http.ResponseWriter) *template.Template {
	tmp, err := template.ParseFiles("templates/" + filename)
	if err != nil {
		http.Error(res, "can't parse html template", http.StatusInternalServerError)
	}
	return tmp
}

func renderTemplatesParseGlob(res http.ResponseWriter) *template.Template {
	tmp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		http.Error(res, "can't parse html template", http.StatusInternalServerError)
	}

	return tmp
}

func (h *handlerStruct) home(res http.ResponseWriter, req *http.Request) {
	rows, err := h.dbPool.Query(context.Background(), "SELECT * FROM movies;")
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(res, "No movies", http.StatusNotFound)
		}
		http.Error(res, "Something went wrong", http.StatusInternalServerError)
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var movies films
	for rows.Next() {
		var film film
		err := rows.Scan(&film.ID, &film.Name, &film.Director, &film.Rating)
		// nil deference means accessing a pointer of a variable which is not inintialized yet
		if err != nil {
			http.Error(res, "oh no!", http.StatusInternalServerError)
			fmt.Println(err.Error())
			break
		}
		movies = append(movies, film)
	}

	dataToPasedtoTemplate := map[string]any{"movies": movies, "title": "Home"}

	tpl := renderTemplatesParseGlob(res)
	tpl.ExecuteTemplate(res, "home.html", dataToPasedtoTemplate)
	// tpl.Execute(res, nil) // it will execute the template according to the route and filename
}

func (hand *handlerStruct) getMovieWithId(res http.ResponseWriter, req *http.Request) {
	movieId := req.PathValue("id")
	result := hand.dbPool.QueryRow(context.Background(), "SELECT id, name, director, rating FROM movies WHERE id = $1;", movieId)

	var data film
	err := result.Scan(&data.ID, &data.Name, &data.Director, &data.Rating)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	// res.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(res).Encode(data)
	// fmt.Println(req.Header.Get("Accept"))

	tpl := renderTemplatesParseGlob(res)
	tpl.ExecuteTemplate(res, "movie.html", data)
}

func (h *handlerStruct) postMovie(res http.ResponseWriter, req *http.Request) {
	rating, err := strconv.ParseFloat(req.FormValue("rating"), 64)
	if err != nil {
		http.Error(res, "rating is not valid", http.StatusBadRequest)
		return
	}

	movie := newMovie(req.FormValue("movie"), req.FormValue("director"), rating)
	quryString := "INSERT INTO movies (name, director, rating) VALUES ($1, $2, $3)"
	_, err = h.dbPool.Exec(context.Background(), quryString, movie.Name, movie.Director, movie.Rating)
	// fmt.Println(a) // it shows number of items inserted
	if err != nil {
		http.Error(res, "could not add to db", http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, "/", http.StatusSeeOther) // use proper redirect status codes to make it work
}
