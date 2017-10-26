package lexer

import (
	"strings"

	"github.com/rmerry/tncbot/ircclient/tokens"
)

var commands = map[string]tokens.Type{
	"PING":    tokens.PING,
	"PRIVMSG": tokens.PRIVMSG}

// Token represents and IRC message token
type Token struct {
	Type  tokens.Type
	Value string
}

// Lexer holds properties and functions related to the lexical analysis of an
// IRC message
type Lexer struct {
	input    string
	location int
}

// New returns an instance of type Lexer
func New(input string) *Lexer {
	return &Lexer{
		input: input}
}

// NextToken returns
func (l *Lexer) NextToken() *Token {
	for i := l.location; ; i++ {
		if l.input[i] != ' ' {
			l.location = i
			break
		}
	}

	// tokens beginning with `:'
	if l.input[l.location] == ':' {
		// it's a prefix if we are at the start of the line
		if l.location == 0 {
			l.location = l.location + 1
			return &Token{tokens.PREFIX, l.getPrefix()}
		}

		l.location = l.location + 1
		return &Token{tokens.STRING, l.getString()}
	}

	if l.input[l.location] == '#' || l.input[l.location] == '&' {
		return &Token{tokens.CHANNEL, l.getChannelName()}
	}

	if token, ok := l.getCommand(); ok {
		return &Token{Type: token}
	}

	return &Token{Type: tokens.EOF}
}

func (l *Lexer) getChannelName() string {
	var channelName string
	for i := l.location; ; i++ {
		if i == len(l.input) || l.input[i] == ' ' {
			channelName = l.input[l.location:i]
			l.location = i
			break
		}
	}

	return channelName
}

func (l *Lexer) getPrefix() string {
	var prefix string
	for i := l.location; ; i++ {
		if i == len(l.input) || l.input[i] == ' ' {
			prefix = l.input[l.location:i]
			l.location = i
			break
		}
	}

	return prefix
}

func (l *Lexer) getCommand() (tokens.Type, bool) {
	var testStr string
	for i := l.location; ; i++ {
		if i == len(l.input) || l.input[i] == ' ' {
			testStr = l.input[l.location:i]
			break
		}
	}

	// check that string is an IRC command
	token, ok := commands[testStr]
	if !ok {
		return tokens.INVALID, false
	}

	l.location = l.location + len(testStr)

	return token, true
}

// FIXME: THIS IS BAD, FIX IT
func (l *Lexer) getString() string {
	var str string
	for i := l.location; ; i++ {
		if i == len(l.input) {
			str = l.input[l.location:i]
			l.location = i
			break
		}
	}

	return strings.TrimRight(str, "\r\n")
}
