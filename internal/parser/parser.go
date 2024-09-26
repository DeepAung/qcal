package parser

import (
	"fmt"

	"github.com/DeepAung/qcal/internal/ast"
	"github.com/DeepAung/qcal/internal/lexer"
	"github.com/DeepAung/qcal/internal/token"
)

// order of precedence
const (
	_ int = iota
	LOWEST
	OR      // or
	AND     // and
	COMPARE // ==, !=, <, <=, >, >=
	SUM     // +, -
	PRODUCT // *, /, %
	PREFIX  // -5, !true
	POWER   // ^
	POSTFIX // 5!
	CALL    // myFunc()
)

var precedences = map[token.TokenType]int{
	token.OR:       OR,
	token.AND:      AND,
	token.EQ:       COMPARE,
	token.NOT_EQ:   COMPARE,
	token.LT:       COMPARE,
	token.GT:       COMPARE,
	token.LT_EQ:    COMPARE,
	token.GT_EQ:    COMPARE,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.PERCENT:  PRODUCT,
	// PREFIXS
	token.CARET: POWER,
	// POSTFIXS
	token.LPAREN: CALL,
}

type (
	prefixParseFn  func() ast.Expression
	postfixParseFn func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.TokenType]prefixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: make([]string, 0)}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() (*ast.Program, []string) {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program, p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	if p.curToken.Type == token.IDENT && p.peekToken.Type == token.ASSIGN {
		return p.parseLetStatement()
	}

	return p.parseExpressionStatement()
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	return &ast.ExpressionStatement{Token: p.curToken}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	return nil
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type != t {
		msg := fmt.Sprintf("expect next token to be %s, got %s instead", t, p.peekToken.Type)
		p.errors = append(p.errors, msg)
		return false
	}

	p.nextToken()
	return true
}
