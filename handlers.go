package main

import (
	"context"
	"html/template"
	"net/http"
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

func (hand *handler) home(res http.ResponseWriter, req *http.Request) {
	// res.Write([]byte("hello"))
	tpl := renderTemplatesParseGlob(res)
	// tpl.Execute(res, nil) // it will execute the template according to the route and filename
	tpl.ExecuteTemplate(res, "home.html", nil)
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

func (hand *handler) yupp(res http.ResponseWriter, req *http.Request) {
	// templ := renderTemplates("index.html", res)
	// templ.Execute(res, nil)

	tmpl2 := renderTemplatesParseGlob(res)
	tmpl2.Execute(res, nil)
}
