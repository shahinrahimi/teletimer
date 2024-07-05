package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateDirecryIfNotExist(directory string) error {
	// check if directory exist
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// try creating a new directory
		if err := os.Mkdir(directory, 0755); err != nil {
			log.Println("Error creating directory for database", err)
			return err
		} else {
			log.Println("New directory for database created!")
		}
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func AreSlicesEqual(slice1, slice2 []int64) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, v := range slice1 {
		if v != slice2[i] {
			return false
		}
	}
	return true
}

func ParseDuration(durationStr string) (time.Duration, error) {
	unit := defaultUnit
	if len(durationStr) > 1 {
		unit = durationStr[len(durationStr)-1:]
		durationStr = durationStr[:len(durationStr)-1]
	}

	num, err := strconv.Atoi(durationStr)
	if err != nil {
		return 0, err
	}

	switch unit {
	case "s":
		return time.Duration(num) * time.Second, nil
	case "min":
		return time.Duration(num) * time.Minute, nil
	case "h":
		return time.Duration(num) * time.Hour, nil
	case "d":
		return time.Duration(num) * time.Hour * 24, nil
	default:
		return 0, fmt.Errorf("invalid unit")
	}
}
