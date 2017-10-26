package main

import (
	"github.com/rmerry/tncbot/bot"
	"github.com/rmerry/tncbot/config"
)

func main() {
	c, err := config.Parse()
	if err != nil {

	}

	bot, err := bot.New(&bot.NewOptions{
		Password: c.Password,
		Channel:  c.Channel,
		Ident:    c.Ident,
		Nick:     c.Nickname,
		Port:     c.Port,
		Server:   c.Server})
	if err != nil {

	}

	err = bot.Connect()
	if err != nil {

	}
	// bot.Start()
}
