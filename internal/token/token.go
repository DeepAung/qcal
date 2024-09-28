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
	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER" // e.g. "123", "112.", ".20", "122.02"

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
	LT_EQ  TokenType = "<="
	GT_EQ  TokenType = ">="
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	ARROW  TokenType = "=>"
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	// Keywords
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	OR     TokenType = "OR"
	AND    TokenType = "AND"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	RETURN TokenType = "RETURN"
)

var keywords = map[string]TokenType{
	"true":   TRUE,
	"false":  FALSE,
	"or":     OR,
	"and":    AND,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(literal string) TokenType {
	if tokenType, ok := keywords[literal]; ok {
		return tokenType
	}
	return IDENT
}
