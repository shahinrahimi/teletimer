package main

import (
	"log"
	"os"

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
