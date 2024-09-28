package parser

import (
	"fmt"
	"strconv"

	"github.com/DeepAung/qcal/internal/ast"
	"github.com/DeepAung/qcal/internal/lexer"
	"github.com/DeepAung/qcal/internal/token"
)

// order of precedence
const (
	_ int = iota
	LOWEST
	OR  // or
	AND // and
	EQUALS
	COMPARE  // ==, !=, <, <=, >, >=
	SUM      // +, -
	PRODUCT  // *, /, %
	PREFIX   // -5, !true
	EXPONENT // ^
	POSTFIX  // 5!
	CALL     // myFunc()

	ARROW_FUNCTION
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
	token.CARET:    EXPONENT,
	token.BANG:     POSTFIX,
	token.LPAREN:   CALL,

	token.ARROW: ARROW_FUNCTION,
}

var associativity = map[int]string{
	LOWEST:   "left",
	OR:       "left",
	AND:      "left",
	EQUALS:   "left",
	COMPARE:  "left",
	SUM:      "left",
	PRODUCT:  "left",
	PREFIX:   "left",
	EXPONENT: "right", // !important
	POSTFIX:  "left",
	CALL:     "left",
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken   token.Token
	peekToken  token.Token
	peek2Token token.Token
	peek3Token token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: make([]string, 0),

		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.NUMBER, p.parseNumber)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpressionOrFunctionLiteral)

	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.PLUS, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.registerInfix(token.BANG, p.parsePostfixExpression)

	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.CARET, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.ARROW, p.parseInfixFunctionLiteral)

	p.nextToken()
	p.nextToken()
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

	if p.curToken.Type == token.RETURN {
		return p.parseReturnStatement()
	}

	if p.curToken.Type == token.SEMICOLON {
		return nil
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf(
			"no prefix parse function for %s %q found",
			p.curToken.Type,
			p.curToken.Literal,
		)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLON &&
		(precedence < p.peekPrecedence() || (precedence == p.peekPrecedence() && associativity[precedence] == "right")) {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
		if _, ok := leftExp.(ast.FunctionLiteral); ok {
			return leftExp
		}
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == token.TRUE}
}

func (p *Parser) parseNumber() ast.Expression {
	lit := &ast.NumberLiteral{Token: p.curToken}

	number, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = number
	return lit
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	if p.curToken.Type != token.LPAREN {
		exp.Condition = p.parseExpression(LOWEST)
	} else {
		p.nextToken()
		exp.Condition = p.parseExpression(LOWEST)

		if !p.expectPeek(token.RPAREN) {
			return nil
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekToken.Type != token.ELSE {
		return exp
	}
	p.nextToken()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Alternative = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: make([]ast.Statement, 0),
	}

	p.nextToken()

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		if stmt := p.parseStatement(); stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

/*
Function Literal
() => {}
(i) => { i }
(i) => i
(i, j) => {}
i => {} // infix => operator ???

Grouped Expression
()
(1 + i)
(i + 1)
*/

/*
()
() =>
(a)
(a) =>
(a, b) =>
*/
func (p *Parser) parseGroupedExpressionOrFunctionLiteral() ast.Expression {
	if p.peekToken.Type == token.RPAREN {
		if p.peek2Token.Type == token.ARROW {
			return p.parseFunctionLiteral()
		}

		p.errors = append(p.errors, "invalid () grouped expression")
		return nil
	}

	if p.peekToken.Type == token.IDENT {
		if p.peek2Token.Type == token.COMMA {
			return p.parseFunctionLiteral()
		}

		if p.peek2Token.Type == token.RPAREN {
			if p.peek3Token.Type == token.ARROW {
				return p.parseFunctionLiteral()
			}

			return p.parseGroupedExpression()
		}

		return p.parseGroupedExpression()
	}

	return p.parseGroupedExpression()
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	parameters := p.parseFunctionParameters()

	p.nextToken()
	tok := p.curToken // the arrow token

	if p.peekToken.Type == token.LBRACE {
		p.nextToken()
		body := p.parseBlockStatement()
		return &ast.NormalFunctionLiteral{Token: tok, Parameters: parameters, Body: body}
	}

	p.nextToken()
	body := p.parseExpression(LOWEST)
	return &ast.ConciseFunctionLiteral{Token: tok, Parameters: parameters, Body: body}
}

func (p *Parser) parseInfixFunctionLiteral(left ast.Expression) ast.Expression {
	param, ok := left.(*ast.Identifier)
	if !ok {
		return nil
	}
	parameters := []*ast.Identifier{param}

	tok := p.curToken

	if p.peekToken.Type == token.LBRACE {
		p.nextToken()
		body := p.parseBlockStatement()
		return &ast.NormalFunctionLiteral{Token: tok, Parameters: parameters, Body: body}
	}

	p.nextToken()
	body := p.parseExpression(LOWEST)
	return &ast.ConciseFunctionLiteral{Token: tok, Parameters: parameters, Body: body}
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	idents := make([]*ast.Identifier, 0)

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return idents
	}
	p.nextToken()

	idents = append(idents, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		idents = append(idents, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return idents
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	if p.peekToken.Type == token.RPAREN {
		p.errors = append(p.errors, "invalid () grouped expression")
		return nil
	}

	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	funcIdent, ok := left.(*ast.Identifier)
	if !ok {
		p.errors = append(p.errors, fmt.Sprintf("cannot call a function of %q", left.String()))
		return nil
	}

	exp := &ast.CallExpression{Token: p.curToken, Function: funcIdent}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Type == end {
		p.nextToken()
		return nil
	}
	p.nextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}

// Helper functions ----------------------------------------------------------------- //

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.peek2Token
	p.peek2Token = p.peek3Token
	p.peek3Token = p.l.NextToken()
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

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) registerPrefix(t token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t token.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}
