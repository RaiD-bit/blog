package main

import (
	"database/sql"
	"html/template"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Article struct {
	ID      int           `json:"id"`
	Title   string        `json:"title"`
	Content template.HTML `json:"content"`
}

var (
	router *chi.Mux
	db     *sql.DB
)
