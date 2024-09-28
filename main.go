package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DeepAung/qcal/internal/evaluator"
	"github.com/DeepAung/qcal/internal/lexer"
	"github.com/DeepAung/qcal/internal/object"
	"github.com/DeepAung/qcal/internal/parser"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()

	for {
		fmt.Print("input math expression: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "" {
			break
		}

		l := lexer.New(input)
		p := parser.New(l)
		program, errors := p.ParseProgram()
		if errors != nil && len(errors) > 0 {
			fmt.Println("error:")
			for _, msg := range errors {
				fmt.Println("\t- ", msg)
			}
		}

		obj := evaluator.Eval(program, env)
		if obj == nil {
			continue
		}

		fmt.Println(obj.Inspect())
	}
}
