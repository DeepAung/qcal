package ast

import (
	"strings"

	"github.com/DeepAung/qcal/internal/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}

	return p.Statements[0].TokenLiteral()
}

func (p *Program) String() string {
	var sb strings.Builder

	for _, s := range p.Statements {
		sb.WriteString(s.String())
	}

	return sb.String()
}

// LetStatement `<identifier | Name> = <expression | Value>`
type LetStatement struct {
	Token token.Token // The identifier token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var sb strings.Builder

	sb.WriteString(ls.Name.String())
	sb.WriteString(" = ")

	if ls.Value != nil {
		sb.WriteString(ls.Value.String())
	}

	sb.WriteString(";")

	return sb.String()
}

// ExpressionStatement
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// Identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// BooleanLiteral
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string       { return b.Token.Literal }

// NumberLiteral
type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (il *NumberLiteral) expressionNode()      {}
func (il *NumberLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *NumberLiteral) String() string       { return il.Token.Literal }

// PrefixExpression `<prefix | Operator><expression | Right>`
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. `!` from `!true`
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(pe.Operator)
	sb.WriteString(pe.Right.String())
	sb.WriteString(")")

	return sb.String()
}

// PostfixExpression `<expression | Left><postfix | Operator>`
type PostfixExpression struct {
	Token    token.Token // The postfix token, e.g. `!` from `5!`
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(pe.Left.String())
	sb.WriteString(pe.Operator)
	sb.WriteString(")")

	return sb.String()
}

// InfixExpression `<expression | Left> <operator | Operator> <expression | Right>`
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. `+`
	Operator string
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ie.Left.String())
	sb.WriteString(" " + ie.Operator + " ")
	sb.WriteString(ie.Right.String())
	sb.WriteString(")")

	return sb.String()
}

// IfExpression `if (<condition>) { <consequence> } else { <alternative> }`
type IfExpression struct {
	Token       token.Token // the `if` token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var sb strings.Builder

	sb.WriteString("if (")
	sb.WriteString(ie.Condition.String())
	sb.WriteString(") { ")
	sb.WriteString(ie.Consequence.String())
	sb.WriteString(" }")
	if ie.Alternative != nil {
		sb.WriteString(" else { ")
		sb.WriteString(ie.Alternative.String())
		sb.WriteString(" }")
	}

	return sb.String()
}

type FunctionLiteral struct {
	Token      token.Token // the `(` token
	Parameters []*Identifier
	Body       Expression
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var sb strings.Builder

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	sb.WriteString("(")
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") => ")
	sb.WriteString(fl.Body.String())

	return sb.String()
}

type CallExpression struct {
	Token     token.Token // the `(` token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var sb strings.Builder

	args := []string{}
	for _, p := range ce.Arguments {
		args = append(args, p.String())
	}

	sb.WriteString(ce.Function.String())
	sb.WriteString("(")
	sb.WriteString(strings.Join(args, ", "))
	sb.WriteString(")")

	return sb.String()
}
