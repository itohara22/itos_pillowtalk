package main

import "net/http"

func router(h *handler) *http.ServeMux {
	routerMux := http.NewServeMux()

	routerMux.HandleFunc("GET /", h.home)

	// routerMux.HandleFunc("GET /{id}", func(res http.ResponseWriter, req *http.Request) {
	// param := req.PathValue("id")
	routerMux.HandleFunc("GET /{id}", h.nina)

	routerMux.HandleFunc("GET /yup", h.yupp)

	routerMux.HandleFunc("POST /movie", h.postMovie)

	return routerMux
}
