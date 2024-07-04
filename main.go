package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file", err)
	}
	apiKey := os.Getenv("TELEGRAM_BOT_API_KEY")
	if apiKey == "" {
		log.Panic("Telegram bot apiKey not found.")
	}
	store, err := NewSqliteStore()
	if err != nil {
		log.Panic("The database store not found", err)
	}
	if err := store.Init(); err != nil {
		log.Panic("Cant initilized store", err)
	}
	bot, err := NewTelegramBot(store, apiKey)
	if err != nil {
		log.Panic("Cant create a new bot", err)
	}
	go bot.Init()
	user, err := NewUser(13654, "Test", "default")
	if err != nil {
		log.Panic("The user can not created", err)
	}

	if err := bot.store.CreateUser(*user); err != nil {
		log.Panic("Cant create a user", err)
	}

	users, err := bot.store.GetUsers()
	if err != nil {
		log.Panic("cant get users", err)
	}
	fmt.Println(len(users))

}
