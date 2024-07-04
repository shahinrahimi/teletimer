package main

import (
	"log"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type TelegramBot struct {
	store Storage
	bot   *tele.Bot
}

func NewTelegramBot(store Storage, apiKey string) (*TelegramBot, error) {
	pref := tele.Settings{
		Token:  apiKey,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &TelegramBot{
		store: store,
		bot:   b,
	}, nil
}

func (b *TelegramBot) Init() {
	b.bot.Use(middleware.AutoRespond())
	b.bot.Use(AutoResponder)
	usersOnly := b.bot.Group()
	usersOnly.Use(middleware.Whitelist())
	usersOnly.Handle("/addalert", func(c tele.Context) error {
		args := strings.Fields(c.Text())
		if len(args) < 3 {
			return c.Send("Usage: /addalert <label> <number><unit>")
		}
		return nil
		// label := args[1]
		// durationStr := args[2]
		// userID := c.Sender().ID

	})
	usersOnly.Handle("/deletealert", func(c tele.Context) error {
		args := strings.Fields(c.Text())
		if len(args) < 3 {
			return c.Send("Usage: /deletealert <id>")
		}
		return nil
	})
	usersOnly.Handle("/viewalerts", func(c tele.Context) error {
		args := strings.Fields(c.Text())
		if len(args) < 3 {
			return c.Send("Usage: /viewalerts <id>")
		}
		return nil
	})

	b.bot.Start()
}

func AutoResponder(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Callback() != nil {
			defer c.Respond()
		}
		return next(c) // continue execution chain
	}
}

func Logger(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		var (
			user = c.Sender()
			text = c.Text()
		)
		log.Println(user, " wrote ", text)
		return next(c)
	}
}
