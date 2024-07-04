package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteStoreTesting() (*SqliteStore, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return &SqliteStore{
		db,
	}, nil
}
