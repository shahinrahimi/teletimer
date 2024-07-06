package main

import tele "gopkg.in/telebot.v3"

type HandleFunc func(c tele.Context) error

func MakeHandleFunc(f HandleFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if err := f(c); err != nil {
			return err
		}
		return nil
	}
}
