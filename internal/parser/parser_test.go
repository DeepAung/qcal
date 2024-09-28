package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/DeepAung/qcal/internal/ast"
	"github.com/DeepAung/qcal/internal/lexer"
	"github.com/DeepAung/qcal/internal/token"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"y = z;", "y", "z"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.LetStatement{})

		stmt := program.Statements[0].(*ast.LetStatement)

		testIdentifier(t, stmt.Name, tt.expectedIdentifier)
		testLiteralExpression(t, stmt.Value, tt.expectedValue)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.ReturnStatement{})

		stmt := program.Statements[0].(*ast.ReturnStatement)

		if stmt.TokenLiteral() != "return" {
			t.Fatalf(
				"invalid statement token literal, expect='return', got=%q",
				stmt.TokenLiteral(),
			)
		}

		testLiteralExpression(t, stmt.Value, tt.expectedValue)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	expectedIdentifier := "foobar"

	l := lexer.New(input)
	p := New(l)
	program, errors := p.ParseProgram()
	checkParserErrors(t, errors)
	testProgramStatement(t, program, &ast.ExpressionStatement{})

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	testIdentifier(t, stmt.Expression, expectedIdentifier)
}

func TestNumberLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		isError  bool
		expected float64
	}{
		{"5", false, 5},
		{"5.", false, 5},
		{".5", false, 0.5},
		{"5.5", false, 5.5},
		{".", true, 0},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		if tt.isError {
			if errors == nil || len(errors) == 0 {
				t.Fatalf("should error, got no error")
			}
			continue
		}

		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.ExpressionStatement{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		testNumberLiteral(t, stmt.Expression, tt.expected)
	}
}

func TestBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.ExpressionStatement{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		testBooleanLiteral(t, stmt.Expression, tt.expect)
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"+15;", "+", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"+foobar;", "+", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.ExpressionStatement{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf(
				"invalid stmt.Expression type, expect=*ast.PrefixExpression, got=%T",
				stmt.Expression,
			)
		}
		if exp.TokenLiteral() != tt.operator {
			t.Fatalf(
				"invalid exp.TokenLiteral(), expect=%q, got=%q",
				tt.operator, exp.TokenLiteral(),
			)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("invalid exp.Operator, expect=%q, got=%q", tt.operator, exp.Operator)
		}
		testLiteralExpression(t, exp.Right, tt.value)
	}
}

func TestPostfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"5!;", "!", 5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.ExpressionStatement{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		exp, ok := stmt.Expression.(*ast.PostfixExpression)
		if !ok {
			t.Fatalf(
				"invalid stmt.Expression type, expect=*ast.PostfixExpression, got=%T",
				stmt.Expression,
			)
		}
		if exp.TokenLiteral() != tt.operator {
			t.Fatalf(
				"invalid exp.TokenLiteral(), expect=%q, got=%q",
				tt.operator, exp.TokenLiteral(),
			)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("invalid exp.Operator, expect=%q, got=%q", tt.operator, exp.Operator)
		}
		testLiteralExpression(t, exp.Left, tt.value)
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 % 5;", 5, "%", 5},
		{"5 ^ 5;", 5, "^", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar % barfoo;", "foobar", "%", "barfoo"},
		{"foobar ^ barfoo;", "foobar", "^", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar >= barfoo;", "foobar", ">=", "barfoo"},
		{"foobar <= barfoo;", "foobar", "<=", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.ExpressionStatement{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		// simple
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"-a!",
			"(-(a!))",
		},
		{
			"a!!",
			"((a!)!)",
		},
		{
			"a!^2",
			"((a!) ^ 2)",
		},
		{
			"a^2!",
			"(a ^ (2!))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a % b / c * d",
			"(((a % b) / c) * d)",
		},
		{
			"a ^ b ^ c",
			"(a ^ (b ^ c))",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a * b - c",
			"((a * b) - c)",
		},
		{
			"a + b * c + d / f - g",
			"(((a + (b * c)) + (d / f)) - g)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4 >= 5 != 5 <= 5",
			"((((((5 > 4) == 3) < 4) >= 5) != 5) <= 5)",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		// with boolean
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		// with parentheses
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(2 ^ 2) ^ 2",
			"((2 ^ 2) ^ 2)",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		// with call expression
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		checkParserErrors(t, errors)

		got := program.String()
		if got != tt.expect {
			t.Errorf("expect=%q, got=%q", tt.expect, got)
		}
	}
}

func TestIfExpression(t *testing.T) {
	inputs := []string{"if x < y { x }", "if (x < y) { x }"}

	for _, input := range inputs {
		l := lexer.New(input)
		p := New(l)
		program, errors := p.ParseProgram()
		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.ExpressionStatement{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf(
				"invalid stmt.Expression type, expect=*ast.IfExpression, got=%T",
				stmt.Expression,
			)
		}

		testInfixExpression(t, exp.Condition, "x", "<", "y")

		if len(exp.Consequence.Statements) != 1 {
			t.Fatalf(
				"invalid exp.Consequence.Statements length, expect=1, got=%d",
				len(exp.Consequence.Statements),
			)
		}
		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"invalid consequence statement type, expect=*ast.ExpressionStatement, got=%T",
				exp.Consequence.Statements[0],
			)
		}
		testIdentifier(t, consequence.Expression, "x")

		if exp.Alternative != nil {
			t.Fatalf("invalid exp.Alternative, expect=<nil>, got=%+v", exp.Alternative)
		}
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if x < y { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program, errors := p.ParseProgram()
	checkParserErrors(t, errors)
	testProgramStatement(t, program, &ast.ExpressionStatement{})

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf(
			"invalid stmt.Expression type, expect=*ast.IfExpression, got=%T",
			stmt.Expression,
		)
	}

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf(
			"invalid exp.Consequence.Statements length, expect=1, got=%d",
			len(exp.Consequence.Statements),
		)
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"invalid consequence statement type, expect=*ast.ExpressionStatement, got=%T",
			exp.Consequence.Statements[0],
		)
	}
	testIdentifier(t, consequence.Expression, "x")

	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf(
			"invalid exp.Alternative.Statements length, expect=1, got=%d",
			len(exp.Alternative.Statements),
		)
	}
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"invalid alternative statement type, expect=*ast.ExpressionStatement, got=%T",
			exp.Alternative.Statements[0],
		)
	}
	testIdentifier(t, alternative.Expression, "y")
}

func TestNormalFunctionLiteral(t *testing.T) {
	tests := []struct {
		input      string
		parameters []string
		body       string
	}{
		{"() => {}", []string{}, ""},
		{"(i) => {}", []string{"i"}, ""},
		{"(i, j, k) => {}", []string{"i", "j", "k"}, ""},
		{"i => {}", []string{"i"}, ""},
		{"i => { i + 1 }", []string{"i"}, "(i + 1)"},
		{"i => { i + 1; };", []string{"i"}, "(i + 1)"},
		{"i => { i + 1; i }", []string{"i"}, "(i + 1)i"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()
		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.ExpressionStatement{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		function, ok := stmt.Expression.(*ast.NormalFunctionLiteral)
		if !ok {
			t.Fatalf(
				"invalid stmt.Expression type, expect=*ast.NormalFunctionLiteral, got=%T",
				stmt.Expression,
			)
		}

		if function.TokenLiteral() != string(token.ARROW) {
			t.Fatalf(
				"invalid function token literal, expect=%q, got=%q",
				token.ARROW, function.TokenLiteral(),
			)
		}

		for i, param := range function.Parameters {
			if param.String() != tt.parameters[i] {
				t.Fatalf(
					"invalid function parameters[%d], expect=%q, got=%q",
					i, tt.parameters[i], param.String(),
				)
			}
		}

		if function.Body.String() != tt.body {
			t.Fatalf("invalid function body, expect=%q, got=%q", tt.body, function.Body.String())
		}
	}
}

func TestConciseFunctionLiteral(t *testing.T) {
	tests := []struct {
		input      string
		isError    bool
		parameters []string
		body       string
	}{
		{"i => ", true, nil, ""},
		{"i => i + 1", false, []string{"i"}, "(i + 1)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, errors := p.ParseProgram()

		if tt.isError {
			if errors == nil || len(errors) == 0 {
				t.Fatalf("should error, got no error")
			}
			continue
		}

		checkParserErrors(t, errors)
		testProgramStatement(t, program, &ast.ExpressionStatement{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		function, ok := stmt.Expression.(*ast.ConciseFunctionLiteral)
		if !ok {
			t.Fatalf(
				"invalid stmt.Expression type, expect=*ast.ConciseFunctionLiteral, got=%T",
				stmt.Expression,
			)
		}

		if function.TokenLiteral() != string(token.ARROW) {
			t.Fatalf(
				"invalid function token literal, expect=%q, got=%q",
				token.ARROW, function.TokenLiteral(),
			)
		}

		for i, param := range function.Parameters {
			if param.String() != tt.parameters[i] {
				t.Fatalf(
					"invalid function parameters[%d], expect=%q, got=%q",
					i, tt.parameters[i], param.String(),
				)
			}
		}

		if function.Body.String() != tt.body {
			t.Fatalf("invalid function body, expect=%q, got=%q", tt.body, function.Body.String())
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program, errors := p.ParseProgram()
	checkParserErrors(t, errors)
	testProgramStatement(t, program, &ast.ExpressionStatement{})

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf(
			"invalid stmt.Expression type, expect=*ast.CallExpression, got=%T",
			stmt.Expression,
		)
	}

	testIdentifier(t, exp.Function, "add")

	if len(exp.Arguments) != 3 {
		t.Fatalf("invalid exp.Arguments length, expect=3, got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

// ------------------------------------------------------------------ //

func checkParserErrors(t *testing.T, errors []string) {
	t.Helper()

	if errors != nil && len(errors) > 0 {
		t.Error("parseProgram failed:\n")
		for _, msg := range errors {
			t.Errorf("- %s\n", msg)
		}
	}
}

func testProgramStatement(t *testing.T, program *ast.Program, expectedType interface{}) {
	t.Helper()

	if len(program.Statements) != 1 {
		t.Fatalf("invalid program.Statements length, expect=1, got=%d", len(program.Statements))
	}

	expectedTypeName := reflect.TypeOf(expectedType).Name()
	if expectedTypeName != reflect.TypeOf(program.Statements[0]).Name() {
		t.Fatalf(
			"invalid statement type, expect=%q, got=%T",
			expectedTypeName, program.Statements[0],
		)
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) {
	t.Helper()

	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("invalid expression type, expect=*ast.Identifier, got=%T", exp)
	}

	if ident.Value != val {
		t.Fatalf("invalid identifier value, expect=%s, got=%s", val, ident.Value)
	}
	if ident.TokenLiteral() != val {
		t.Fatalf("invalid identifier literal, expect=%s, got=%s", val, ident.TokenLiteral())
	}
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) {
	t.Helper()

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("invalid exp type, expect=*ast.InfixExpression, got=%T", exp)
	}

	if opExp.Operator != operator {
		t.Fatalf("invalid opExp.Operator, expect=%q, got=%q", operator, opExp.Operator)
	}
	if opExp.TokenLiteral() != operator {
		t.Fatalf(
			"invalid opExp.TokenLiteral(), expect=%q, got=%q",
			operator, opExp.TokenLiteral(),
		)
	}

	testLiteralExpression(t, opExp.Left, left)
	testLiteralExpression(t, opExp.Right, right)
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expectedVal interface{}) {
	t.Helper()

	switch v := expectedVal.(type) {
	case int:
		testNumberLiteral(t, exp, float64(v))
	case float64:
		testNumberLiteral(t, exp, v)
	case bool:
		testBooleanLiteral(t, exp, v)
	case string:
		testIdentifier(t, exp, v)
	default:
		t.Fatalf("type of exp is not handled, exp=%T, expectedVal=%T", exp, expectedVal)
	}
}

func testNumberLiteral(t *testing.T, exp ast.Expression, val float64) {
	t.Helper()

	li, ok := exp.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("invalid expression type, expect=*ast.NumberLiteral, got=%T", exp)
	}

	if li.Value != val {
		t.Fatalf("invalid number literal value, expect=%v, got=%v", val, li.Value)
	}
	literalVal, err := strconv.ParseFloat(li.TokenLiteral(), 64)
	if err != nil {
		t.Fatalf("cannot parse number literal to float64: %v", err)
	}
	if literalVal != val {
		t.Fatalf("invalid number literal, expect=%f, got=%f", val, literalVal)
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, val bool) {
	t.Helper()

	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("invalid expression type, expect=*ast.BooleanLiteral, got=%T", exp)
	}

	if bo.Value != val {
		t.Fatalf("invalid boolean literal value, expect=%t, got=%t", val, bo.Value)
	}
	if bo.TokenLiteral() != fmt.Sprint(val) {
		t.Fatalf("invalid boolean literal, expect=%s, got=%s", fmt.Sprint(val), bo.TokenLiteral())
	}
}
