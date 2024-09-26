package calculator

import (
	"log"
	"testing"
)

type testGroup struct {
	name  string
	tests []testObj
}

type testObj struct {
	input    string
	expected float64
}

func TestCalculate(t *testing.T) {
	testGroups := []testGroup{
		{
			name: "Simple operators",
			tests: []testObj{
				{"1+5", 6},
				{"1-9", -8},
				{"9-1", 8},
				{"4*8", 32},
				{"4/8", 0.5},
				{"2^3", 8},
			},
		},
		{
			name: "Prefix",
			tests: []testObj{
				{"-1+8", 7},
				{"-4*8", -32},
			},
		},
		{
			name: "Precedence",
			tests: []testObj{
				{"2^0+3*5/(2+4-1)", 4},
			},
		},
		{
			name: "Associativity",
			tests: []testObj{
				{"2^2^3", 256},
			},
		},
		{
			name: "Build-in variables",
			tests: []testObj{
				{"e", 2},
				{"pi", 3},
			},
		},
		{
			name: "Build-in functions",
			tests: []testObj{
				{"log(1)", 0},
				{"log(e)", 1},
			},
		},
		{
			name: "Variables",
			tests: []testObj{
				{"x = 1; x+5", 6},
				{"x = 1; x+1; x+3;", 4},
				{"x = 1; log(1)", 0},
			},
		},
		{
			name: "Functions",
			tests: []testObj{
				{"f = x => x^2; f(3)", 9},
				{"f = (x, y) => i + 2; f(2)", 4},
			},
		},
	}

	calculator := NewCalculator()

	for _, tt := range testGroups {
		t.Run(tt.name, func(t *testing.T) {
			for _, testObj := range tt.tests {
				testCalculate(t, calculator, testObj)
			}
		})
	}
}

func testCalculate(t *testing.T, calculator *calculator, obj testObj) {
	result, err := calculator.Calculate([]byte(obj.input))
	if err != nil {
		t.Fatalf("input %q, client.Calculate failed: %v", obj.input, err)
	}
	if result != obj.expected {
		t.Fatalf("input %q, expected: %f, got: %f", obj.input, obj.expected, result)
	}

	log.Printf("calculate %q got %f", obj.input, result)
}
