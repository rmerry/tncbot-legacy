package lexer

import (
	"testing"

	"github.com/rmerry/tncbot/ircclient/tokens"
)

func TestNextToken(t *testing.T) {
	type tokenTest struct {
		testName      string
		input         string
		expectedToken Token
	}
	testCases := []tokenTest{
		tokenTest{"prefix token test", ":niven.freenode.net", Token{tokens.PREFIX, "niven.freenode.net"}},
		tokenTest{"PRIVMSG token test", "PRIVMSG", Token{tokens.PRIVMSG, ""}},
		tokenTest{"PING token test", "PING", Token{tokens.PING, ""}},
		tokenTest{"CHANNEL token test (#)", "#TheNorthCode", Token{tokens.CHANNEL, "#TheNorthCode"}},
		tokenTest{"STRING token test (#)", " :this is a test string", Token{tokens.STRING, "this is a test string"}}}

	for _, testCase := range testCases {
		lexer := New(testCase.input)
		token := lexer.NextToken()

		if token.Type != testCase.expectedToken.Type {
			t.Errorf("EXPECTED token of type `%s', GOT token of type `%s'\n", testCase.expectedToken.Type, token.Type)
		}
		if token.Value != testCase.expectedToken.Value {
			t.Errorf("EXPECTED token with value `%s', GOT token with value `%s'\n", testCase.expectedToken.Value, token.Value)
		}
	}
}
