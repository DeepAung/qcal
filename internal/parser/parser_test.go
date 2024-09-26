package parser

import (
	"fmt"
	"testing"

	"github.com/DeepAung/qcal/internal/ast"
	"github.com/DeepAung/qcal/internal/lexer"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		if errors != nil && len(errors) > 0 {
			t.Error("parseProgram failed:\n")
			for _, msg := range errors {
				t.Errorf("- %s\n", msg)
			}
		}

		if len(program.Statements) != 1 {
			t.Fatalf(
				"invalid program.Statements length, expect=1, got=%d",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf(
				"invalid statement type, expect=*ast.LetStatement, got=%T",
				program.Statements[0],
			)
		}

		if stmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf(
				"invalid identifier value, expect=%q, got=%q",
				tt.expectedIdentifier, stmt.Name.Value,
			)
		}
		if stmt.Name.TokenLiteral() != tt.expectedIdentifier {
			t.Fatalf(
				"invalid identifier token literal, expect=%q, got=%q",
				tt.expectedIdentifier, stmt.Name.TokenLiteral(),
			)
		}

		testLiteralExpression(t, stmt.Value, tt.expectedValue)

	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expectedVal interface{}) {
	t.Helper()

	switch v := expectedVal.(type) {
	case float64:
		testNumberLiteral(t, exp, v)
	case bool:
		testBooleanLiteral(t, exp, v)
	default:
		t.Fatalf("type of exp is not handled, exp=%T, expectedVal=%T", exp, expectedVal)
	}
}

func testNumberLiteral(t *testing.T, exp ast.Expression, val float64) {
	t.Helper()

	il, ok := exp.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("invalid exp type, expect=*ast.NumberLiteral, got=%T", exp)
	}

	if il.Value != val {
		t.Fatalf("invalid il.Value, expect=%v, got=%v", val, il.Value)
	}
	if il.TokenLiteral() != fmt.Sprint(val) {
		t.Fatalf("invalid il.TokenLiteral(), expect=%q, got=%q", fmt.Sprint(val), il.TokenLiteral())
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, val bool) {
	t.Helper()

	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("invalid exp type, expect=*ast.BooleanLiteral, got=%T", exp)
	}

	if bo.Value != val {
		t.Fatalf("invalid bo.Value, expect=%t, got=%t", val, bo.Value)
	}
	if bo.TokenLiteral() != fmt.Sprint(val) {
		t.Fatalf("invalid bo.TokenLiteral(), expect=%s, got=%s", fmt.Sprint(val), bo.TokenLiteral())
	}
}
