package bot

import (
	"time"

	"github.com/rmerry/tncbot/bot/actions/chucknorris"
	"github.com/rmerry/tncbot/bot/actions/johnmoss"
	"github.com/rmerry/tncbot/bot/actions/linkdescribe"
	"github.com/rmerry/tncbot/ircclient"
)

const gecos = "TheNorthCodeBot (tncbot) v0.0.1a"

// Action represents the actions that can be undertaken by the bot
type Action struct {
	name     string
	function func(string) (string, bool)
}

var actions = []*Action{
	&Action{"johnmoss", johnmoss.Execute},
	&Action{"chucknorris", chucknorris.Execute},
}

// NewOptions ...
type NewOptions struct {
	Nick     string
	Channel  string
	Port     int
	Server   string
	Ident    string
	Password string
}

// New ..
func New(opts *NewOptions) (*Bot, error) {
	bot := &Bot{
		ident:    opts.Ident,
		password: opts.Password,
		channel:  opts.Channel,
		nick:     opts.Nick,
		server:   opts.Server,
		port:     opts.Port}

	return bot, nil
}

// Bot ..
type Bot struct {
	ircclient *ircclient.IRCClient
	nick      string
	channel   string
	port      int
	server    string
	ident     string
	password  string
}

// Connect ...
func (b *Bot) Connect() error {
	client := ircclient.New(&ircclient.NewIRCClientOptions{
		Gecos:    gecos,
		Password: b.password,
		Ident:    b.ident,
		Nick:     b.nick,
		Port:     b.port,
		Server:   b.server})

	err := client.Connect()
	if err != nil {
		panic(err)
	}
	b.ircclient = client

	time.Sleep(time.Second * 10)
	client.Join(b.channel)

	for {
		msg := <-b.ircclient.Messages
		if linkDescription, ok := linkdescribe.Execute(msg.Value); ok {
			b.ircclient.SendMessage(&ircclient.IRCMessage{Value: linkDescription})
		} else {
			for _, a := range actions {
				if str, ok := a.function(msg.Value); ok {
					b.ircclient.SendMessage(&ircclient.IRCMessage{Value: str})
					break
				}
			}
		}
	}

	return nil
}

// Disconnect ...
func (b *Bot) Disconnect() error {
	return nil
}
