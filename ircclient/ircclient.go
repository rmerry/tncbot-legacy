package ircclient

import (
	"bufio"
	"fmt"
	"net"

	"github.com/rmerry/tncbot/ircclient/lexer"
	"github.com/rmerry/tncbot/ircclient/tokens"
)

// IRCClient ...
type IRCClient struct {
	password  string
	channel   string
	conn      net.Conn
	connected bool
	ident     string
	gecos     string
	nick      string
	port      int
	server    string
	inbox     chan string
	outbox    chan string
	Messages  chan *IRCMessage
}

// IRCMessage ...
type IRCMessage struct {
	Source string
	Target string
	Value  string
}

// NewIRCClientOptions ...
type NewIRCClientOptions struct {
	Gecos    string
	Ident    string
	Nick     string
	Port     int
	Server   string
	Password string
}

// New ...
func New(opts *NewIRCClientOptions) *IRCClient {
	client := &IRCClient{
		password: opts.Password,
		gecos:    opts.Gecos,
		ident:    opts.Ident,
		nick:     opts.Nick,
		port:     opts.Port,
		server:   opts.Server,
		inbox:    make(chan string, 100),
		outbox:   make(chan string, 100),
		Messages: make(chan *IRCMessage, 100)}

	return client
}

// Connect ...
func (c *IRCClient) Connect() error {
	fmt.Println("dialing")
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.server, c.port))
	if err != nil {
		return err
	}

	c.conn = conn
	c.connected = true

	go c.startWriter()
	go c.startListener()

	c.outbox <- fmt.Sprintf("NICK %s\r\n", c.nick)
	c.outbox <- fmt.Sprintf("USER %s * 8 :%s\r\n", c.ident, c.gecos)

	if c.password != "" {
		c.outbox <- fmt.Sprintf("PRIVMSG NickServ :identify %s %s\r\n", c.nick, c.password)
	}

	return nil
}

// SendMessage ...
func (c *IRCClient) SendMessage(m *IRCMessage) {
	c.outbox <- fmt.Sprintf("PRIVMSG %s :%s\r\n", c.channel, m.Value)
}

// Join is used for joining an IRC channel
func (c *IRCClient) Join(channel string) error {
	if c.connected {
		c.outbox <- fmt.Sprintf("JOIN %s\r\n", channel)
		c.channel = channel
		return nil
	}

	return fmt.Errorf("must be connected in order to join channel")
}

func (c *IRCClient) startListener() {
	bufReader := bufio.NewReader(c.conn)

	for {
		msg, err := bufReader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		fmt.Print(msg)
		m := c.parse(msg)
		if m != nil {
			c.Messages <- m
		}
	}
}

func (c *IRCClient) startWriter() {
	bufWriter := bufio.NewWriter(c.conn)

	for {
		select {
		case msg := <-c.outbox:
			_, err := bufWriter.WriteString(msg)
			if err != nil {
				panic(err)
			}

			bufWriter.Flush()
		}
	}
}

func (c *IRCClient) parse(msg string) *IRCMessage {
	lex := lexer.New(msg)

	m := &IRCMessage{}
	t := lex.NextToken()

	if t.Type == tokens.EOF {
		return nil
	}

	if t.Type == tokens.PREFIX {
		m.Source = t.Value
		t = lex.NextToken()
	}

	switch t.Type {
	case tokens.EOF:
		return nil

	case tokens.PRIVMSG:
		for {
			t = lex.NextToken()
			if t.Type == tokens.STRING {
				m.Value = t.Value
				return m
			} else if t.Type == tokens.CHANNEL {
				m.Target = t.Value
			} else {
				return nil
			}
		}

	case tokens.PING:
		c.outbox <- "PONG"
	}

	return nil
}
