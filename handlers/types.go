package handlers

import (
	"html/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HandlerStruct struct {
	DbPool *pgxpool.Pool
	T      *template.Template
}
