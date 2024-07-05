package main

import (
	"fmt"
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

	adminIDs, err := b.store.GetAdminIDs()
	if err != nil {
		log.Println("Error finding usersIDs for admins", err)
	}

	userIDs, err := b.store.GetUserIDs()
	if err != nil {
		log.Println("Error finding usersIDs for users", err)
	}

	b.bot.Handle("/echo", b.HandleEcho())
	b.bot.Handle("/addalert", b.HandleAddAlert())

	adminsOnly := b.bot.Group()
	adminsOnly.Use(middleware.Whitelist(adminIDs...))
	adminsOnly.Handle("/helloadmin", b.HandleHelloAdmin())
	adminsOnly.Handle("/test", b.HandleTest())
	usersOnly := b.bot.Group()
	usersOnly.Use(middleware.Whitelist(userIDs...))
	usersOnly.Handle("/test", b.HandleTest())
	// usersOnly.Handle("/addalert", b.HandleAddAlert())
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

func (b *TelegramBot) HandleEcho() tele.HandlerFunc {
	return func(c tele.Context) error {
		args := strings.Fields(c.Text())
		if len(args) < 2 {
			return c.Send("Usage: /echo <somthing>")
		}
		somthing := args[1]
		return c.Send(somthing)
	}
}

func (b *TelegramBot) HandleHelloAdmin() tele.HandlerFunc {
	return func(c tele.Context) error {
		return c.Send("Hi Admin")
	}
}
func (b *TelegramBot) HandleTest() tele.HandlerFunc {
	return func(c tele.Context) error {
		return c.Send("1\n2\n3\n testing")
	}
}
func (b *TelegramBot) HandleAddAlert() tele.HandlerFunc {
	return func(c tele.Context) error {
		args := strings.Fields(c.Text())
		if len(args) < 3 {
			return c.Send("Usage: /addalert <label> <number><unit>")
		}
		label := args[1]
		durationStr := args[2]
		userID := c.Sender().ID
		duration, err := ParseDuration(durationStr)
		if err != nil {
			return c.Send("Invalid duration format.")
		}
		count, err2 := b.store.GetAlertsCountByUserID(userID)
		if err2 != true {
			log.Println("Does not able to count the alerts by userID")
		}
		if count >= maxAlerts {
			return c.Send("You can only create up to 5 alerts.")
		}
		triggerAt := time.Now().Add(duration)
		newAlert := NewAlert(userID, label, triggerAt)
		if err := b.store.CreateAlert(*newAlert); err != nil {
			return err
		}
		go b.ScheduleAlert(userID, label, triggerAt)
		return c.Send(fmt.Sprintf("Alert created: %s in %s", label, durationStr))

	}
}

func (b *TelegramBot) SendAlert(userID int64, label string) {
	inlineKeys := [][]tele.InlineButton{
		{
			tele.InlineButton{Unique: "snooze", Text: "Snooze"},
			tele.InlineButton{Unique: "dismiss", Text: "Dismiss"},
		},
	}
	b.bot.Send(tele.ChatID(userID), fmt.Sprintf("Alert: %s", label), &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
}

func (b *TelegramBot) ScheduleAlert(userID int64, label string, triggerAt time.Time) {
	time.Sleep(time.Until(triggerAt))
	b.SendAlert(userID, label)
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
