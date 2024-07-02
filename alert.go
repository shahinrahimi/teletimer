package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Alert struct {
	ID     string `json:"id"`
	UserID int64  `json:"user_id"`
}

func NewAlert(userID int64) *Alert {
	return &Alert{
		ID:     fmt.Sprint("TA" + strconv.Itoa(rand.Int())),
		UserID: userID,
	}
}
