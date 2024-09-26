package lexer

import (
	"log"
	"testing"

	"github.com/DeepAung/qcal/internal/token"
)

func TestNextToken(t *testing.T) {
	input := `
x = 123;
y = 123.
y = .20;
y = .;
y = 123.20
x = y + y - x * y / x ^ y % y
f = (a, b) => a + b;
f(e, pi)
	`
	expects := []token.Token{
		{Type: token.IDENT, Literal: "x"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.NUMBER, Literal: "123"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.NUMBER, Literal: "123."},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.NUMBER, Literal: ".20"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.ILLEGAL, Literal: "."},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.NUMBER, Literal: "123.20"},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.ASTERISK, Literal: "*"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.SLASH, Literal: "/"},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.CARET, Literal: "^"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.PERCENT, Literal: "%"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.IDENT, Literal: "f"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "a"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "b"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.ARROW, Literal: "=>"},
		{Type: token.IDENT, Literal: "a"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.IDENT, Literal: "b"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IDENT, Literal: "f"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.CONSTANT, Literal: "e"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.CONSTANT, Literal: "pi"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.EOF, Literal: ""},
	}

	l := New(input)
	for i, expect := range expects {
		tok := l.NextToken()
		if tok.Type != expect.Type {
			log.Fatalf(
				"expects[%d] - invalid token type, expect=%q, got=%q",
				i, expect.Type, tok.Type,
			)
		}
		if tok.Literal != expect.Literal {
			log.Fatalf(
				"expects[%d] - invalid token literal, expect=%q, got=%q",
				i, expect.Literal, tok.Literal,
			)
		}
	}
}
