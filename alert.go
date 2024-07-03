package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Alert struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Lable     string    `json:"label"`
	TriggerAt time.Time `json:"trigger_at"`
}

func GetCreateTableAlertsQuery() string {
	return `CREATE TABLE IF NOT EXISTS alerts (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		label TEXT,
		trigger_at DATETIME,
	);`
}

func NewAlert(userID int64) *Alert {
	return &Alert{
		ID:     fmt.Sprint("TA" + strconv.Itoa(rand.Int())),
		UserID: userID,
	}
}
