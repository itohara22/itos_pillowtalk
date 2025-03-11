package handlers

import (
	"context"
	"net/http"
	"sushi/utils"

	"github.com/jackc/pgx/v5"
)

func (h *HandlerStruct) RegisterUser(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	id := 0
	query := "SELECT * FROM users WHERE username = $1;"
	rows := h.DbPool.QueryRow(context.Background(), query, username)
	err := rows.Scan(id)

	if err == pgx.ErrNoRows {
		hasedP, err := utils.HashPassword(password)
		if err != nil {
			http.Error(res, "something went wrong", http.StatusInternalServerError)
			return
		}
		qur := "insert into users (username, password) values($1,$2)"
		_, err = h.DbPool.Exec(context.Background(), qur, username, hasedP)
		if err != nil {
			http.Error(res, "could not add to db", http.StatusInternalServerError)
			return
		}

	} else if err == nil {
		http.Error(res, "User already exists", http.StatusBadRequest)
	} else {
		http.Error(res, "something went wrong", http.StatusInternalServerError)
	}
}
