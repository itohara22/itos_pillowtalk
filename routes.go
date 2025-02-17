package main

import (
	"net/http"
	"sushi/handlers"
)

func router(h *handlers.HandlerStruct) *http.ServeMux {
	dir := http.Dir("./static")
	fs := http.FileServer(dir)
	routerMux := http.NewServeMux()

	routerMux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	routerMux.HandleFunc("GET /", h.Home)

	// routerMux.HandleFunc("GET /{id}", func(res http.ResponseWriter, req *http.Request) {
	// param := req.PathValue("id")
	routerMux.HandleFunc("GET /film/{id}", h.GetMovieWithId)

	routerMux.HandleFunc("POST /movie", h.PostMovie)

	return routerMux
}
