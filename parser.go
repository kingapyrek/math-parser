package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

// Interface for numberNode and operatorNode
type node interface {
	eval() int
}

type numberNode struct {
	value int
}

type operatorNode struct {
	left      node
	right     node
	operation byte
}

// evaluate node with number
func (node *numberNode) eval() int {
	return node.value
}

// evaluate node with operator
func (node *operatorNode) eval() int {
	leftResult := node.left.eval()
	rightResult := node.right.eval()

	switch node.operation {
	case '+':
		return leftResult + rightResult
	case '-':
		return leftResult - rightResult
	case '*':
		return leftResult * rightResult
	case '/':
		return leftResult / rightResult
	default:
		panic("not known operation")
	}
}

type parser struct {
	expression string
	pos        int
}

func (p *parser) currentChar() byte {
	if p.pos >= len(p.expression) {
		return 0
	} else {
		return p.expression[p.pos]
	}
}

func (p *parser) moveCursor() {
	p.pos++
	//ignore whitespace
	if p.currentChar() == ' ' {
		p.pos++
	}
}

func (p *parser) parsePrimary() (node, error) {
	curr := p.currentChar()
	if curr == '(' {
		//move to expression inside ()
		p.moveCursor()
		node, err := p.parseAddSubOp()
		if err != nil {
			return nil, err
		}
		if p.currentChar() != ')' {
			return nil, fmt.Errorf("invalid character, expected ')' after '(")
		}
		//move cursor after expression and char ')'
		p.moveCursor()
		return node, nil
	} else if curr >= '0' && curr <= '9' {
		//handle digit
		value, err := strconv.Atoi(string(curr))
		if err != nil {
			return nil, fmt.Errorf("conversion char to integer failed")
		}
		// Validation whether number is one digit and not fractional
		if p.pos+1 < len(p.expression) {
			nextChar := p.expression[p.pos+1]
			if nextChar >= '0' && nextChar <= '9' {
				return nil, fmt.Errorf("only single-digit integers are allowed")
			} else if nextChar == '.' {
				return nil, fmt.Errorf("fractional numbers are not allowed")
			}
		}
		p.moveCursor()
		return &numberNode{value: value}, nil
	} else {
		return nil, fmt.Errorf("unexpected char in expression")
	}
}

func (p *parser) parseMulDiv() (node, error) {
	node, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for p.currentChar() == '*' || p.currentChar() == '/' {
		operation := p.currentChar()
		p.moveCursor()
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		node = &operatorNode{
			left:      node,
			operation: operation,
			right:     right,
		}
	}
	return node, nil
}

func (p *parser) parseAddSubOp() (node, error) {
	node, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}
	for p.currentChar() == '+' || p.currentChar() == '-' {
		operation := p.currentChar()
		p.moveCursor()
		right, err := p.parseMulDiv()
		if err != nil {
			return nil, err
		}
		node = &operatorNode{
			left:      node,
			operation: operation,
			right:     right,
		}
	}
	return node, nil
}

func (p *parser) parse() (node, error) {
	return p.parseAddSubOp()
}

func main() {
	file, err := os.Open("equations.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		equation := scanner.Text()

		// Increment the WaitGroup counter
		wg.Add(1)

		// Create goroutine to parralel parsing and evaluating equations
		go func(eq string) {
			defer wg.Done()

			parser := &parser{expression: eq, pos: 0}
			tree, err := parser.parse()
			if err != nil {
				fmt.Printf("Error parsing expression %v: %v\n", eq, err)
				return
			}
			result := tree.eval()

			fmt.Printf("%s = %d\n", eq, result)
		}(equation)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
