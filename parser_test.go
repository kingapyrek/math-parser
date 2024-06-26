package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Define expected error messages
var (
	errUnexpectedCharacter  = errors.New("unexpected character in expression")
	errFractionalNumber     = errors.New("fractional numbers are not allowed")
	errMultiDigitNumber     = errors.New("only single-digit integers are allowed")
	errUnmatchedParenthesis = errors.New("invalid character, expected ')' after '('")
	errNotOpenedParanthesis = errors.New("unexpected ')' without opening '(")
)

// TestParseAddSub tests the parseAddSub function with various inputs
func TestParseAddSub(t *testing.T) {
	tests := []struct {
		expression  string
		expected    int
		expectedErr error
	}{
		{"1+2", 3, nil},
		{"1-2", -1, nil},
		{"1+2-3", 0, nil},
		{"2*3+1", 7, nil},
		{"2+3*4", 14, nil},
		{"(2+3)*4", 20, nil},
		{"(2+3)*(1+1)", 10, nil},
		{"1+", 0, errUnexpectedCharacter},
		{"1++2", 0, errUnexpectedCharacter},
		{"1+2.0", 0, errFractionalNumber},
		{"12+3", 0, errMultiDigitNumber},
		{"(1 + 5 * 3", 0, errUnmatchedParenthesis},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			p := &parser{expression: test.expression, pos: 0}
			node, err := p.parseAddSubOp()
			if test.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error for expression: %s, but got none", test.expression)
				} else if err.Error() != test.expectedErr.Error() {
					t.Errorf("Expected error %v for expression: %s, but got %v", test.expectedErr, test.expression, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for expression: %s, error: %v", test.expression, err)
					return
				}
				result := node.eval()

				if result != test.expected {
					t.Errorf("Expected result %d for expression: %s, but got %d", test.expected, test.expression, result)
				}
			}
		})
	}
}

// TestParseMulDiv tests the parseMulDiv function with various inputs
func TestParseMulDiv(t *testing.T) {
	tests := []struct {
		expression     string
		expectedResult int
		expectedErr    error
	}{
		{"2*3", 6, nil},
		{"6/3", 2, nil},
		{"2*3/4", 1, nil},
		{"2/3*4", 0, nil},
		{"(2*3)/4", 1, nil},
		{"2*(3/4)", 0, nil},
		{"2*3.0", 0, errFractionalNumber},
		{"12*3", 0, errMultiDigitNumber},
		{"(1 * 5 / 3", 0, errUnmatchedParenthesis},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			p := &parser{expression: test.expression, pos: 0}
			node, err := p.parseMulDiv()
			if test.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error for expression: %s, but got none", test.expression)
				} else if err.Error() != test.expectedErr.Error() {
					t.Errorf("Expected error %v for expression: %s, but got %v", test.expectedErr, test.expression, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for expression: %s, error: %v", test.expression, err)
					return
				}
				result := node.eval()

				if result != test.expectedResult {
					t.Errorf("Expected result %d for expression: %s, but got %d", test.expectedResult, test.expression, result)
				}
			}
		})
	}
}

// TestParsePrimary tests the parsePrimary function with various inputs
func TestParsePrimary(t *testing.T) {
	tests := []struct {
		expression     string
		expectedResult int
		expectedErr    error
	}{
		{"1", 1, nil},
		{"(1)", 1, nil},
		{"(1+2)", 3, nil},
		{"(1+(2*3))", 7, nil},
		{"1.2+2", 0, errFractionalNumber},
		{"12+3", 0, errMultiDigitNumber},
		{"(1 + 5 * 3", 0, errUnmatchedParenthesis},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			p := &parser{expression: test.expression, pos: 0}
			node, err := p.parsePrimary()
			if test.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error for expression: %s, but got none", test.expression)
				} else if err.Error() != test.expectedErr.Error() {
					t.Errorf("Expected error %v for expression: %s, but got %v", test.expectedErr, test.expression, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for expression: %s, error: %v", test.expression, err)
					return
				}
				result := node.eval()

				if result != test.expectedResult {
					t.Errorf("Expected result %d for expression: %s, but got %d", test.expectedResult, test.expression, result)
				}
			}
		})
	}
}

// TestParse tests the parse function with various inputs
func TestParse(t *testing.T) {
	tests := []struct {
		expression     string
		expectedResult int
		expectedErr    error
	}{
		{"1+2", 3, nil},
		{"1-2", -1, nil},
		{"1+2-3", 0, nil},
		{"2*3+1", 7, nil},
		{"2+3*4", 14, nil},
		{"(2+3)*4", 20, nil},
		{"(2+3)*(1+1)", 10, nil},
		{"1+", 0, errUnexpectedCharacter},
		{"1++2", 0, errUnexpectedCharacter},
		{"1+2.0", 0, errFractionalNumber},
		{"12+3", 0, errMultiDigitNumber},
		{"(1 + 5 * 3", 0, errUnmatchedParenthesis},
		{"1+2)", 0, errNotOpenedParanthesis},
		{"1a + 4", 0, errUnexpectedCharacter},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			p := &parser{expression: test.expression, pos: 0}
			node, err := p.parse()
			if test.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error for expression: %s, but got none", test.expression)
				} else if err.Error() != test.expectedErr.Error() {
					t.Errorf("Expected error %v for expression: %s, but got %v", test.expectedErr, test.expression, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for expression: %s, error: %v", test.expression, err)
					return
				}
				result := node.eval()

				if result != test.expectedResult {
					t.Errorf("Expected result %d for expression: %s, but got %d", test.expectedResult, test.expression, result)
				}
			}
		})
	}
}

func TestOperatorNodeEval(t *testing.T) {
	tests := []struct {
		name         string
		leftValue    int
		rightValue   int
		operation    byte
		expected     int
		expectPanic  bool
		errorMessage string
	}{
		{"Addition", 2, 3, '+', 5, false, ""},
		{"Subtraction", 5, 3, '-', 2, false, ""},
		{"Multiplication", 4, 6, '*', 24, false, ""},
		{"Division", 10, 2, '/', 5, false, ""},
		{"DivisionByZero", 10, 0, '/', 0, true, "runtime error: integer divide by zero"},
		{"UnknownOperation", 2, 3, '%', 0, true, "not known operation"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			leftNode := &numberNode{value: test.leftValue}
			rightNode := &numberNode{value: test.rightValue}
			opNode := &operatorNode{left: leftNode, right: rightNode, operation: test.operation}

			if test.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic, but did not panic")
					} else {
						if !strings.Contains(fmt.Sprint(r), test.errorMessage) {
							t.Errorf("Expected panic message '%v', but got '%v'", test.errorMessage, r)
						}
					}
				}()
			}

			result := opNode.eval()

			if !test.expectPanic && result != test.expected {
				t.Errorf("Expected result %d, but got %d", test.expected, result)
			}
		})
	}
}

func TestNumberNodeEval(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"PositiveNumber", 5, 5},
		{"NegativeNumber", -3, -3},
		{"Zero", 0, 0},
		{"LargeNumber", 1000000, 1000000},
		{"MinInt", -2147483648, -2147483648}, // testing with minimum integer value
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			numNode := &numberNode{value: test.value}

			result := numNode.eval()

			if result != test.expected {
				t.Errorf("Expected result %d, but got %d", test.expected, result)
			}
		})
	}
}

func TestCurrentChar(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		pos        int
		expected   byte
	}{
		{"ValidPosition", "abc", 1, 'b'},
		{"EndOfString", "def", 3, 0},
		{"EmptyString", "", 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &parser{expression: test.expression, pos: test.pos}

			result := parser.currentChar()

			if result != test.expected {
				t.Errorf("Expected character '%c', but got '%c'", test.expected, result)
			}
		})
	}
}

func TestMoveCursor(t *testing.T) {
	tests := []struct {
		name         string
		expression   string
		initialPos   int
		finalPos     int
		expectedChar byte
	}{
		{"IncrementPosition", "abc", 0, 1, 'b'},
		{"IgnoreWhitespace", "d ef", 0, 1, 'e'},
		{"MoveToEnd", "ghi", 2, 3, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &parser{expression: test.expression, pos: test.initialPos}

			parser.moveCursor()

			currentChar := parser.currentChar()
			if currentChar != test.expectedChar {
				t.Errorf("Expected current character '%c', but got '%c'", test.expectedChar, currentChar)
			}
		})
	}
}
