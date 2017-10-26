// Package tokens ..
package tokens

// Type represents an IRC message token type
type Type int

const (
	// PREFIX ..
	PREFIX Type = iota
	// CHANNEL ...
	CHANNEL
	// EOF ..
	EOF
	// PRIVMSG ..
	PRIVMSG
	// PING ..
	PING
	// INVALID ...
	INVALID
	// STRING ...
	STRING
)

// Token represents a token of an IRC message
type Token struct {
	Type  Type
	Value string
}
