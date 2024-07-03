package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
}

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore() (*SqliteStore, error) {
	if err := CreateDirecryIfNotExist("database"); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", "./database/mydb.db")
	if err != nil {
		return nil, err
	}
	log.Println("DB Connected!")

	return &SqliteStore{
		db,
	}, nil
}

func (s *SqliteStore) Init() error {
	if _, err := s.db.Exec(GetCreateTableUsersQuery()); err != nil {
		log.Println("Cant create users table", err)
		return err
	}
	if _, err := s.db.Exec(GetCreateTableAlertsQuery()); err != nil {
		log.Println("Cant create alerts table", err)
		return err
	}
	return nil
}
