package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	GetAdminIDs() ([]int64, error)
	GetUserIDs() ([]int64, error)
	GetUser(id string) (*User, error)
	GetUsers() ([]User, error)
	GetUserByUserID(userID int64) (*User, error)
	CreateUser(user User) error
	UpdateUser(id string, u User) error
	DeleteUser(id string) error
	GetUsersCount() (int, bool)

	GetAlert(id string) (*Alert, error)
	GetAlerts() ([]Alert, error)
	GetAlertsByUserID(userID int64) ([]Alert, error)
	GetAlertsCountByUserID(userID int64) (int, bool)
	CreateAlert(a Alert) error
	UpdateAlert(id string, a Alert) error
	DeleteAlert(id string) error
	DeleteAlertsByUserID(userID int64) error
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
func (s *SqliteStore) GetAdminIDs() ([]int64, error) {
	rows, err := s.db.Query("SELECT user_id FROM users WHERE is_admin = ?", true)
	if err != nil {
		return nil, err
	}
	var adminIDs []int64
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID); err != nil {
			return nil, err
		}
		adminIDs = append(adminIDs, user.UserID)
	}
	return adminIDs, nil
}
func (s *SqliteStore) GetUserIDs() ([]int64, error) {
	rows, err := s.db.Query("SELECT user_id FROM users")
	if err != nil {
		return nil, err
	}
	var userIDs []int64
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, user.UserID)
	}
	return userIDs, nil
}
func (s *SqliteStore) GetUsersCount() (int, bool) {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return 0, false
	}
	return count, true
}

func (s *SqliteStore) GetUser(id string) (*User, error) {
	query := GetSelectUserQuery()
	var user User
	feilds := user.ToFeilds()
	err := s.db.QueryRow(query, id).Scan(feilds...)
	if err != nil {
		if err == sql.ErrNoRows {
			// handle user not found
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (s *SqliteStore) GetUsers() ([]User, error) {
	query := GetSelectAllUsersQuery()
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	var users []User
	for rows.Next() {
		var user User
		fields := user.ToFeilds()
		if err := rows.Scan(fields...); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (s *SqliteStore) GetUserByUserID(userID int64) (*User, error) {
	query := GetSelectUserByUserIDQuery()
	var user User
	feilds := user.ToFeilds()
	if err := s.db.QueryRow(query, userID).Scan(feilds...); err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *SqliteStore) CreateUser(user User) error {
	query := GetInsertUserQuery()
	args := user.ToArgs()
	if _, err := s.db.Exec(query, args...); err != nil {
		return err
	}
	return nil
}
func (s *SqliteStore) UpdateUser(id string, user User) error {
	query := GetUpdateUserQuery()
	args := user.ToUpdatedArgs()
	if _, err := s.db.Exec(query, args...); err != nil {
		return err
	}
	return nil
}
func (s *SqliteStore) DeleteUser(id string) error {
	query := GetDeleteUserQuery()
	if _, err := s.db.Exec(query, id); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStore) GetAlert(id string) (*Alert, error) {
	query := GetSelectAlertQuery()
	var alert Alert
	feilds := alert.ToFeilds()
	if err := s.db.QueryRow(query, id).Scan(feilds...); err != nil {
		return nil, err
	}
	return &alert, nil
}
func (s *SqliteStore) GetAlerts() ([]Alert, error) {
	query := GetSelectAlertsQuery()
	var alerts []Alert
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var alert Alert
		feilds := alert.ToFeilds()
		if err := rows.Scan(feilds...); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}
func (s *SqliteStore) GetAlertsByUserID(userID int64) ([]Alert, error) {
	query := GetSelectAlertsByUserIDQuery()
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	var alerts []Alert
	for rows.Next() {
		var alert Alert
		if err := rows.Scan(alert.ToFeilds()...); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}
func (s *SqliteStore) GetAlertsCountByUserID(userID int64) (int, bool) {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM alerts WHERE user_id = ?", userID).Scan(&count); err != nil {
		return 0, false
	}
	return count, true
}
func (s *SqliteStore) CreateAlert(alert Alert) error {
	query := GetInsertAlertQuery()
	if _, err := s.db.Exec(query, alert.ToArgs()...); err != nil {
		return err
	}
	return nil
}
func (s *SqliteStore) UpdateAlert(id string, alert Alert) error {
	query := GetUpdateAlertQuery()
	if _, err := s.db.Exec(query, alert.ToArgs()...); err != nil {
		return err
	}
	return nil
}
func (s *SqliteStore) DeleteAlert(id string) error {
	query := GetDeleteAlertQuery()
	if _, err := s.db.Exec(query, id); err != nil {
		return err
	}
	return nil
}
func (s *SqliteStore) DeleteAlertsByUserID(userID int64) error {
	if _, err := s.db.Exec("DELETE FROM alerts WHERE user_id = ?", userID); err != nil {
		return err
	}
	return nil
}
