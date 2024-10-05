package object

import (
	"fmt"
	"strings"

	"github.com/DeepAung/qcal/internal/ast"
)

type ObjectType string

const (
	NUMBER_OBJ           ObjectType = "NUMBER"
	BOOLEAN_OBJ          ObjectType = "BOOLEAN"
	NULL_OBJ             ObjectType = "NULL"
	ERROR_OBJ            ObjectType = "ERROR"
	LET_VALUE_OBJ        ObjectType = "LET_VAULE"
	RETURN_VALUE_OBJ     ObjectType = "RETURN_VALUE"
	FUNCTION_OBJ         ObjectType = "FUNCTION"
	BUILTIN_FUNCTION_OBJ ObjectType = "BUILTIN_FUNCTION"
	BUILTIN_VALUE_OBJ    ObjectType = "BUILTIN_VALUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Number struct {
	Value float64
}

func (i *Number) Type() ObjectType { return NUMBER_OBJ }
func (i *Number) Inspect() string  { return fmt.Sprint(i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprint(b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type LetValue struct {
	Value Object
}

func (lv *LetValue) Type() ObjectType { return LET_VALUE_OBJ }
func (lv *LetValue) Inspect() string  { return lv.Value.Inspect() }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type NormalFunction struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
}

func (f *NormalFunction) Type() ObjectType { return FUNCTION_OBJ }
func (f *NormalFunction) Inspect() string {
	var sb strings.Builder

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	sb.WriteString("(")
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") => {\n")
	sb.WriteString(f.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

type ConciseFunction struct {
	Parameters []*ast.Identifier
	Body       ast.Expression
}

func (f *ConciseFunction) Type() ObjectType { return FUNCTION_OBJ }
func (f *ConciseFunction) Inspect() string {
	var sb strings.Builder

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	sb.WriteString("(")
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") => ")
	sb.WriteString(f.Body.String())

	return sb.String()
}

type BuiltinFunction func(args ...Object) Object

func (b BuiltinFunction) Type() ObjectType { return BUILTIN_FUNCTION_OBJ }
func (b BuiltinFunction) Inspect() string  { return "builtin function" }
