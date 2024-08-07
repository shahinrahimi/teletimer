package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	IsAdmin   bool      `json:"is_admin"`
}

func GetCreateTableUsersQuery() string {
	return `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL UNIQUE,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		is_admin BOOLEAN
	);`
}

func GetSelectAllUsersQuery() string {
	return `SELECT id, user_id, username, password, created_at, is_admin FROM users`
}
func GetSelectUserQuery() string {
	return `SELECT id, user_id, username, password, created_at, is_admin FROM users WHERE id = ?`
}
func GetSelectUserByUserIDQuery() string {
	return `SELECT id, user_id, username, password, created_at, is_admin FROM users WHERE user_id = ?`
}
func GetInsertUserQuery() string {
	return `INSERT INTO users (id, user_id, username, password, created_at, is_admin) VALUES (?, ?, ?, ?, ?, ?)`
}
func GetUpdateUserQuery() string {
	return `UPDATE users SET user_id = ?, username = ?, password = ?, created_at = ?, is_admin = ? WHERE id = ?`
}
func GetDeleteUserQuery() string {
	return `DELETE FROM users WHERE id = ?`
}

func NewUser(user_id int64, username, password string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        fmt.Sprint("TU" + strconv.Itoa(rand.Int())),
		UserID:    user_id,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		IsAdmin:   false,
	}, nil
}

func NewAdmin(user_id int64, username, password string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        fmt.Sprint("TA" + strconv.Itoa(rand.Int())),
		UserID:    user_id,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		IsAdmin:   true,
	}, nil
}

func (u *User) ToTelegramString() string {
	return fmt.Sprintf("User ID: %d\nUsername: %s\nCreated At: %s",
		u.UserID, u.Username, u.CreatedAt.Format(time.RFC3339))
}

func (u *User) ToArgs() []interface{} {
	return []interface{}{u.ID, u.UserID, u.Username, u.Password, u.CreatedAt, u.IsAdmin}
}
func (u *User) ToUpdatedArgs() []interface{} {
	return []interface{}{u.UserID, u.Username, u.Password, u.CreatedAt, u.IsAdmin, u.ID}
}
func (u *User) ToFeilds() []interface{} {
	return []interface{}{&u.ID, &u.UserID, &u.Username, &u.Password, &u.CreatedAt, &u.IsAdmin}
}
