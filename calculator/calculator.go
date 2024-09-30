package calculator

import (
	"errors"
	"strings"

	"github.com/DeepAung/qcal/internal/evaluator"
	"github.com/DeepAung/qcal/internal/lexer"
	"github.com/DeepAung/qcal/internal/object"
	"github.com/DeepAung/qcal/internal/parser"
)

type Calculator struct {
	env *object.Environment
}

func NewCalculator() *Calculator {
	return &Calculator{
		env: object.NewEnvironment(),
	}
}

func (c *Calculator) Calculate(input string) (object.Object, error) {
	program, errMessages := parser.New(lexer.New(input)).ParseProgram()
	if errMessages != nil && len(errMessages) > 0 {
		if len(errMessages) == 1 {
			return nil, errors.New("ERROR: " + errMessages[0])
		}
		return nil, errors.New("ERROR:\n" + strings.Join(errMessages, "\n"))
	}

	evaluated := evaluator.Eval(program, c.env)
	if evaluated == nil {
		return nil, nil
	}

	if evaluated.Type() == object.ERROR_OBJ {
		return nil, errors.New(evaluated.Inspect())
	}

	return evaluated, nil
}
