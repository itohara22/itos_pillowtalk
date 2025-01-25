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

func (h *handler) home(res http.ResponseWriter, req *http.Request) {
	rows, err := h.dbPool.Query(context.Background(), "SELECT * FROM movies;")
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(res, "No movies", http.StatusNotFound)
		}
		http.Error(res, "Something went wrong", http.StatusInternalServerError)
		fmt.Println(err.Error())
	}

	var movies []movie
	for rows.Next() {
		var film movie
		err := rows.Scan(&film.ID, &film.Name, &film.Director, &film.Rating)
		if err != nil {
			http.Error(res, "oh no!", http.StatusInternalServerError)
			fmt.Println(err.Error())
			break
		}
		movies = append(movies, film)
	}

	tpl := renderTemplatesParseGlob(res)
	tpl.ExecuteTemplate(res, "home.html", movies)
	// tpl.Execute(res, nil) // it will execute the template according to the route and filename
}

func (hand *handler) nina(res http.ResponseWriter, req *http.Request) {
	movieId := req.PathValue("id")
	result := hand.dbPool.QueryRow(context.Background(), "SELECT id, name, director, rating FROM movies WHERE id = $1;", movieId)

	var data movie
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

func (h *handler) postMovie(res http.ResponseWriter, req *http.Request) {
	rating, err := strconv.ParseFloat(req.FormValue("rating"), 64)
	if err != nil {
		http.Error(res, "rating is not valid", http.StatusBadRequest)
	}
	movie := movie{
		Name:     req.FormValue("movie"),
		Director: req.FormValue("director"),
		Rating:   rating,
	}

	quryString := "INSERT INTO movies (name, director, rating) VALUES ($1, $2, $3)"

	_, err = h.dbPool.Exec(context.Background(), quryString, movie.Name, movie.Director, movie.Rating)
	// fmt.Println(a) // it shows number of items inserted
	if err != nil {
		http.Error(res, "could not add to db", http.StatusInternalServerError)
	}

	res.Write([]byte("movie added"))
}

func (hand *handler) yupp(res http.ResponseWriter, req *http.Request) {
	// templ := renderTemplates("index.html", res)
	// templ.Execute(res, nil)

	tmpl2 := renderTemplatesParseGlob(res)
	tmpl2.Execute(res, nil)
}
