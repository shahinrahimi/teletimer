package main

import (
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
	bot.Init()

}
