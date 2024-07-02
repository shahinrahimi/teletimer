package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type User struct {
	Id        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	IsAdmin   bool      `json:"is_admin"`
}

func GetCreateUsersTable() string {
	return `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL UNIQUE,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		is_admin BOOLEAN
	);`
}

func NewUser(user_id int64, username, password string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:        fmt.Sprint("TU" + strconv.Itoa(rand.Int())),
		UserID:    user_id,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		IsAdmin:   false,
	}, nil
}

func NewAdmin(user_id int64, password string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:        fmt.Sprint("AD99999999999"),
		UserID:    user_id,
		Username:  "admin",
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		IsAdmin:   true,
	}, nil
}

func (u *User) toTelegramString() string {
	return fmt.Sprintf("User ID: %d\nUsername: %s\nFistname: %s\nLastname: %s\nCreated At: %s",
		u.UserID, u.Username, u.CreatedAt.Format(time.RFC3339))
}
