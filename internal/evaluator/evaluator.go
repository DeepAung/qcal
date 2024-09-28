package evaluator

import (
	"fmt"
	"math"

	"github.com/DeepAung/qcal/internal/ast"
	"github.com/DeepAung/qcal/internal/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LetStatement:
		if _, ok := builtinValues[node.Name.Value]; ok {
			return newError("cannot assign value to the builtin constant %q", node.Name.Value)
		}
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return &object.LetValue{Value: val}

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}

	case *ast.BooleanLiteral:
		return booleanObject(node.Value)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.NormalFunctionLiteral:
		return &object.NormalFunction{Parameters: node.Parameters, Body: node.Body}

	case *ast.ConciseFunctionLiteral:
		return &object.ConciseFunction{Parameters: node.Parameters, Body: node.Body}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalPostfixExpression(node.Operator, left)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(fn, args)

	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue, *object.Error:
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorPrefixExpression(right)
	case "-":
		return evalMinusOperatorPrefixExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorPrefixExpression(right object.Object) object.Object {
	return booleanObject(!isTruthy(right))
}

func evalMinusOperatorPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

func evalPostfixExpression(operator string, left object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorPostfixExpression(left)
	default:
		return newError("unknown operator: %s%s", left.Type(), operator)
	}
}

// TODO:
func evalBangOperatorPostfixExpression(left object.Object) object.Object {
	if left.Type() != object.NUMBER_OBJ {
		return newError("unknown operator: %s!", left.Type())
	}

	var n int64 = int64(math.Round(left.(*object.Number).Value))
	var result int64 = 1
	for i := int64(2); i <= n; i++ {
		result *= i
	}

	return &object.Number{Value: float64(result)}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalNumberInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Number).Value
	rightValue := right.(*object.Number).Value

	switch operator {
	case "+":
		return &object.Number{Value: leftValue + rightValue}
	case "-":
		return &object.Number{Value: leftValue - rightValue}
	case "*":
		return &object.Number{Value: leftValue * rightValue}
	case "/":
		return &object.Number{Value: leftValue / rightValue}
	case "%":
		return &object.Number{
			Value: float64(int64(math.Round(leftValue)) % int64(math.Round(rightValue))),
		}
	case "^":
		return &object.Number{Value: math.Pow(leftValue, rightValue)}
	case "<":
		return booleanObject(leftValue < rightValue)
	case ">":
		return booleanObject(leftValue > rightValue)
	case "==":
		return booleanObject(leftValue == rightValue)
	case "!=":
		return booleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	switch operator {
	case "==":
		return booleanObject(left == right)
	case "!=":
		return booleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if val, ok := builtinValues[node.Value]; ok {
		return val
	}

	if fn, ok := builtinFuncs[node.Value]; ok {
		return fn
	}

	return newError("identifier not found: %s", node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object

	for _, e := range exps {
		obj := Eval(e, env)
		if isError(obj) {
			return []object.Object{obj}
		}
		results = append(results, obj)
	}

	return results
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.NormalFunction:
		env := createEnv(fn.Parameters, args)
		evaluated := Eval(fn.Body, env)
		return unwrapReturnValue(evaluated)

	case *object.ConciseFunction:
		env := createEnv(fn.Parameters, args)
		evaluated := Eval(fn.Body, env)
		return unwrapReturnValue(evaluated)

	case *object.BuiltinFunction:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func createEnv(params []*ast.Identifier, args []object.Object) *object.Environment {
	env := object.NewEnvironment()
	for i, param := range params {
		env.Set(param.Value, args[i])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnObj, ok := obj.(*object.ReturnValue); ok {
		return returnObj.Value
	}

	return obj
}

// ---------------------------------------------------------------- //

func booleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj == nil {
		return false
	}
	return obj.Type() == object.ERROR_OBJ
}
