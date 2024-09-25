package main

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/DeepAung/calculator-grpc/calculator"
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
				{"2^3^4", 4},
			},
		},
	}

	client, conn := newClient()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for _, tt := range testGroups {
		t.Run(tt.name, func(t *testing.T) {
			for _, testObj := range tt.tests {
				testCalculate(t, client, ctx, testObj)
			}
		})
	}
}

func testCalculate(t *testing.T, client pb.CalculatorClient, ctx context.Context, obj testObj) {
	res, err := client.Calculate(ctx, &pb.Expression{Expression: obj.input})
	if err != nil {
		t.Fatalf("input %q, client.Calculate failed: %v", obj.input, err)
	}
	if res.Result != obj.expected {
		t.Fatalf("input %q, expected: %f, got: %f", obj.input, obj.expected, res.Result)
	}

	log.Printf("calculate %q got %f", obj.input, res.Result)
}
