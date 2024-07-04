package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	GetUser(id string) (*User, error)
	GetUsers() ([]User, error)
	GetUserByUserID(userID int64) (*User, error)
	CreateUser(user User) error
	UpdateUser(id string, u User) error
	DeleteUser(id string) error

	GetAlert(id string) (*Alert, error)
	GetAlerts() ([]Alert, error)
	GetAlertsByUserID(userID int64) ([]Alert, error)
	CreateAlert(a Alert) error
	UpdateAlert(id string, a Alert) error
	DeleteAlert(id string) error
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

func (s *SqliteStore) GetUser(id string) (*User, error) {
	return nil, nil
}
func (s *SqliteStore) GetUsers() ([]User, error) {
	return nil, nil
}
func (s *SqliteStore) GetUserByUserID(userID int64) (*User, error) {
	return nil, nil
}
func (s *SqliteStore) CreateUser(user User) error {
	return nil
}
func (s *SqliteStore) UpdateUser(id string, user User) error {
	return nil
}
func (s *SqliteStore) DeleteUser(id string) error {
	return nil
}

func (s *SqliteStore) GetAlert(id string) (*Alert, error) {
	return nil, nil
}
func (s *SqliteStore) GetAlerts() ([]Alert, error) {
	return nil, nil
}
func (s *SqliteStore) GetAlertsByUserID(userID int64) ([]Alert, error) {
	return nil, nil
}
func (s *SqliteStore) CreateAlert(alert Alert) error {
	return nil
}
func (s *SqliteStore) UpdateAlert(id string, alert Alert) error {
	return nil
}
func (s *SqliteStore) DeleteAlert(id string) error {
	return nil
}
