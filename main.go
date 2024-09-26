package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DeepAung/qcal/calculator"
)

func main() {
	calc := calculator.NewCalculator()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("input math expression: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "" {
			break
		}

		result, err := calc.Calculate([]byte(input))
		if err != nil {
			fmt.Printf("error: %v\n", err)
		} else {
			fmt.Println(result)
		}
	}
}
