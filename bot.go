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

	// adminIDs, err := b.store.GetAdminIDs()
	// if err != nil {
	// 	log.Println("Error finding usersIDs for admins", err)
	// }

	// userIDs, err := b.store.GetUserIDs()
	// if err != nil {
	// 	log.Println("Error finding usersIDs for users", err)
	// }

	// public cammands
	b.bot.Handle("/echo", MakeHandleFunc(b.HandleEcho), Authenticater)
	b.bot.Handle("/start", MakeHandleFunc(b.HandleRegister))
	b.bot.Handle("/deleteme", MakeHandleFunc(b.HandleDeleteMe))

	// usersonly
	usersOnly := b.bot.Group()
	usersOnly.Use(b.RequiredAuthenticated)
	usersOnly.Handle("/addalert", MakeHandleFunc(b.HandleAddAlert))
	usersOnly.Handle("/viewalerts", MakeHandleFunc(b.HandleViewAlerts))
	usersOnly.Handle("/deletealert", MakeHandleFunc(b.HandleDeleteAlert))
	usersOnly.Handle("/updatealert", MakeHandleFunc(b.HandleUpdateAlert))
	usersOnly.Handle(tele.OnCallback, MakeHandleFunc(b.HandleButtons))
	// usersOnly.Use(middleware.Whitelist(userIDs...))

	// adminsonly
	adminsOnly := b.bot.Group()
	adminsOnly.Use(b.RequiredAuthenticated)
	adminsOnly.Handle("/viewusers", MakeHandleFunc(b.HandleViewUsers), b.RequiredAthorization)
	adminsOnly.Handle("/kickuser", MakeHandleFunc(b.HandleKickUser), b.RequiredAthorization)
	adminsOnly.Handle("/banuser", MakeHandleFunc(b.HandleBanUser), b.RequiredAthorization)
	// adminsOnly.Use(middleware.Whitelist(adminIDs...))

	b.bot.Start()
}

// public handlers
func (b *TelegramBot) HandleRegister(c tele.Context) error {
	userID := c.Sender().ID
	if _, err := b.store.GetUserByUserID(userID); err == nil {
		return c.Send("You are already registered!")
	}
	user, err := NewUser(userID, c.Sender().Username, "default_passs")
	if err != nil {
		return c.Send("Failed to create user object.")
	}
	if err := b.store.CreateUser(*user); err != nil {
		return c.Send("Failed to insert user to users.")
	}
	return c.Send("You are registered successfully!")
}
func (b *TelegramBot) HandleDeleteMe(c tele.Context) error {
	return nil
}

// private
func (b *TelegramBot) HandleEcho(c tele.Context) error {
	args := strings.Fields(c.Text())
	if len(args) < 2 {
		return c.Send("Usage: /echo <somthing>")
	}
	somthing := args[1]
	return c.Send(somthing)
}
func (b *TelegramBot) HandleAddAlert(c tele.Context) error {
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
	count, check := b.store.GetAlertsCountByUserID(userID)
	if !check {
		return fmt.Errorf("Faild to count alert by UserID")
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
func (b *TelegramBot) HandleViewAlerts(c tele.Context) error {
	return nil
}
func (b *TelegramBot) HandleDeleteAlert(c tele.Context) error {
	return nil
}
func (b *TelegramBot) HandleUpdateAlert(c tele.Context) error {
	return nil
}
func (b *TelegramBot) HandleButtons(c tele.Context) error {
	return nil
}
func (b *TelegramBot) HandleSnooz() tele.HandlerFunc {
	return func(c tele.Context) error {
		// Handle snooze logic
		data := c.Callback().Data
		fmt.Println("hi there", data)
		alertID := c.Callback().Data
		// if err != nil {
		// 	return c.Send("Invalid alert ID.")
		// }

		alert, err := b.store.GetAlert(alertID)
		if err != nil {
			return c.Send("Alert not found.", alertID)
		}

		duration := alert.TriggerAt.Sub(time.Now()) / snoozeFraction
		newTriggerAt := time.Now().Add(duration)
		alert.TriggerAt = newTriggerAt
		if err := b.store.UpdateAlert(alertID, *alert); err != nil {
			return c.Send("Failed to snooze alert.", err)
		}
		go b.ScheduleAlert(alert.UserID, alert.Lable, newTriggerAt)
		return c.Send(fmt.Sprintf("Alert snoozed for %v", duration))
	}
}

// admin
func (b *TelegramBot) HandleViewUsers(c tele.Context) error {
	return nil
}

func (b *TelegramBot) HandleKickUser(c tele.Context) error {
	return nil
}

func (b *TelegramBot) HandleBanUser(c tele.Context) error {
	return nil
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

func Authenticater(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		var (
			user = c.Sender()
			text = c.Text()
		)
		log.Println(user, " wrote ", text)
		return next(c)
	}
}

func (b *TelegramBot) RequiredAuthenticated(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		var userID = c.Sender().ID
		if _, err := b.store.GetUserByUserID(userID); err != nil {
			c.Send("You are not registered!")
			return err
		}
		return next(c)
	}
}

func (b *TelegramBot) RequiredAthorization(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		var userID = c.Sender().ID
		user, err := b.store.GetUserByUserID(userID)
		if err != nil {
			c.Send("You are not registered!")
			return err
		}
		if !user.IsAdmin {
			c.Send("Permission denied!")
		}
		return next(c)
	}
}
