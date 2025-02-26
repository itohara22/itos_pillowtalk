package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (h *HandlerStruct) RegisterUser(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	// password := req.FormValue("password")

	id := 0
	query := "SELECT * FROM users WHERE username = $1;"
	rows := h.DbPool.QueryRow(context.Background(), query, username)
	err := rows.Scan(id)

	if err == pgx.ErrNoRows {
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("created"))
	} else if err == nil {
		http.Error(res, "something went wrong", http.StatusInternalServerError)
	} else {
		http.Error(res, "User already exists", http.StatusBadRequest)
	}
}
