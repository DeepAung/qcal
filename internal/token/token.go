package token

type Token struct {
	Type    TokenType
	Literal string
}

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + Literals
	IDENT    TokenType = "IDENT"
	NUMBER   TokenType = "NUMBER" // e.g. "123", "112.", ".20", "122.02"
	CONSTANT TokenType = "CONSTANT"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	PERCENT  TokenType = "%"
	CARET    TokenType = "^"
	BANG     TokenType = "!"

	LT     TokenType = "<"
	GT     TokenType = ">"
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	NEWLINE   TokenType = "\n"

	ARROW  TokenType = "=>"
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
	// LBRACKET TokenType = "["
	// RBRACKET TokenType = "]"

	// Keywords
	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"
	IF    TokenType = "IF"
	ELSE  TokenType = "ELSE"
)

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"if":    IF,
	"else":  ELSE,
}

var constants = map[string]struct{}{
	"e":  {},
	"pi": {},
}

func LookupIdent(literal string) TokenType {
	if tokenType, ok := keywords[literal]; ok {
		return tokenType
	}
	if _, ok := constants[literal]; ok {
		return CONSTANT
	}
	return IDENT
}
