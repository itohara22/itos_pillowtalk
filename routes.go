package main

import "net/http"

func router(h *handlerStruct) *http.ServeMux {
	dir := http.Dir("./static")
	fs := http.FileServer(dir)
	routerMux := http.NewServeMux()

	routerMux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	routerMux.HandleFunc("GET /", h.home)

	// routerMux.HandleFunc("GET /{id}", func(res http.ResponseWriter, req *http.Request) {
	// param := req.PathValue("id")
	routerMux.HandleFunc("GET /film/{id}", h.getMovieWithId)

	routerMux.HandleFunc("POST /movie", h.postMovie)

	return routerMux
}
