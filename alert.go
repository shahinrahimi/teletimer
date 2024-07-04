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
		trigger_at DATETIME
	);`
}

func GetSelectAlertsQuery() string {
	return `SELECT id, user_id, label, trigger_at FROM alerts`
}
func GetSelectAlertQuery() string {
	return `SELECT id, user_id, label, trigger_at FROM alerts WHERE id = ?`
}
func GetSelectAlertsByUserIDQuery() string {
	return `SELECT id, user_id, label, trigger_at FROM alerts WHERE user_id = ?`
}
func GetInsertAlertQuery() string {
	return `INSERT INTO alerts (id, user_id, label, trigger_at) VALUES (?,?,?,?)`
}
func GetUpdateAlertQuery() string {
	return `UPDATE alerts SET label = ?, trigger_at = ? WHERE id = ?`
}
func GetDeleteAlertQuery() string {
	return `DELETE FROM alerts WHERE id = ?`
}

func NewAlert(userID int64, label string, triggerAt time.Time) *Alert {
	return &Alert{
		ID:        fmt.Sprint("TA" + strconv.Itoa(rand.Int())),
		UserID:    userID,
		Lable:     label,
		TriggerAt: triggerAt,
	}
}

func (a *Alert) ToArgs() []interface{} {
	return []interface{}{a.ID, a.UserID, a.Lable, a.TriggerAt}
}
func (a *Alert) ToFeilds() []interface{} {
	return []interface{}{&a.ID, &a.UserID, &a.Lable, &a.TriggerAt}
}
