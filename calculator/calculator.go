package calculator

import (
	"errors"
	"fmt"
	"math"

	"github.com/DeepAung/qcal/internal/stack"
)

var ErrInvalidExpression = errors.New("invalid expression")

var precedences = map[byte]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
	'^': 3,
	'(': 4,
}

var associativity = map[byte]byte{
	'+': 'L',
	'-': 'L',
	'*': 'L',
	'/': 'L',
	'^': 'R',
	'(': 'L',
}

type calculator struct{}

func NewCalculator() *calculator {
	return &calculator{}
}

func (s *calculator) Calculate(exp []byte) (float64, error) {
	result, err := parseExpression(exp)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func parseExpression(exp []byte) (float64, error) {
	numStk := stack.NewStack[float64]()
	infixStk := stack.NewStack[byte]()
	var lastCh byte = 0
	var prefix byte = 0

	for i := 0; i < len(exp); i++ {
		cur := exp[i]

		if exp[i] == ' ' || exp[i] == '\t' || exp[i] == '\n' {
			continue
		}

		if isNumber(cur) {
			number, idx := parseNumber(exp, i)
			numStk.Push(number)
			i = idx - 1 // because i++ in for loop

			if prefix != 0 {
				applyPrefix(numStk, &prefix)
			}
		} else if (lastCh == 0 || lastCh == '(') && isPrefix(cur) {
			prefix = cur
		} else if isInfix(cur) {
			for !infixStk.Empty() && precedences[cur] <= precedences[infixStk.Last()] && infixStk.Last() != '(' {
				if err := applyInfix(numStk, infixStk); err != nil {
					return 0, err
				}
			}
			infixStk.Push(cur)
		} else if cur == '(' {
			infixStk.Push(cur)
		} else if cur == ')' {
			for !infixStk.Empty() && infixStk.Last() != '(' {
				if err := applyInfix(numStk, infixStk); err != nil {
					return 0, err
				}
			}

			if infixStk.Empty() || infixStk.Last() != '(' {
				return 0, ErrInvalidExpression
			}
			infixStk.Pop() // pop '('
		} else {
			return 0, fmt.Errorf("invalid token: %c", cur)
		}

		lastCh = cur
	}

	for !infixStk.Empty() {
		if err := applyInfix(numStk, infixStk); err != nil {
			return 0, err
		}
	}

	if numStk.Empty() {
		return 0, ErrInvalidExpression
	}
	return numStk.Last(), nil
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isPrefix(ch byte) bool {
	return ch == '-' || ch == '+'
}

func isInfix(ch byte) bool {
	switch ch {
	case '+', '-', '*', '/', '^':
		return true
	default:
		return false
	}
}

func parseNumber(exp []byte, idx int) (float64, int) {
	var num float64 = 0
	for idx < len(exp) && isNumber(exp[idx]) {
		num = num*10 + float64(exp[idx]-'0')
		idx++
	}
	return num, idx
}

func applyPrefix(numStk *stack.Stack[float64], prefix *byte) error {
	if numStk.Empty() {
		return ErrInvalidExpression
	}

	num := numStk.Pop()
	switch *prefix {
	case '-':
		num = -num
	case '+':
	default:
		return fmt.Errorf("invalid prefix: %d", prefix)
	}

	*prefix = 0
	numStk.Push(num)
	return nil
}

func applyInfix(numStk *stack.Stack[float64], infixStk *stack.Stack[byte]) error {
	if infixStk.Empty() || numStk.Len() < 2 {
		return ErrInvalidExpression
	}

	secondNumber := numStk.Pop()
	firstNumber := numStk.Pop()
	op := infixStk.Pop()

	val, err := evalInfix(firstNumber, op, secondNumber)
	if err != nil {
		return err
	}

	numStk.Push(val)
	return nil
}

func evalInfix(a float64, op byte, b float64) (float64, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		return a / b, nil
	case '^':
		return math.Pow(a, b), nil
	default:
		return 0, fmt.Errorf("invalid operator: %c", op)
	}
}
