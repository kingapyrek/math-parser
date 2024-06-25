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

func (p *parser) parsePrimary() node {
	curr := p.currentChar()
	if curr == '(' {
		//move to expression inside ()
		p.moveCursor()
		node := p.parseAddSubOp()
		if p.currentChar() != ')' {
			panic("Expected ')' after '(")
		}
		//move cursor after expression and char ')'
		p.moveCursor()
		return node
	} else if curr >= '0' && curr <= '9' {
		//handle digit
		p.moveCursor()
		value, _ := strconv.Atoi(string(curr))
		return &numberNode{value: value}
	} else {
		panic("unexpexted char")
	}
}

func (p *parser) parseMulDiv() node {
	node := p.parsePrimary()
	for p.currentChar() == '*' || p.currentChar() == '/' {
		operation := p.currentChar()
		p.moveCursor()
		right := p.parsePrimary()

		node = &operatorNode{
			left:      node,
			operation: operation,
			right:     right,
		}
	}
	return node
}

func (p *parser) parseAddSubOp() node {
	node := p.parseMulDiv()
	for p.currentChar() == '+' || p.currentChar() == '-' {
		operation := p.currentChar()
		p.moveCursor()
		right := p.parseMulDiv()
		node = &operatorNode{
			left:      node,
			operation: operation,
			right:     right,
		}
	}
	return node
}

func (p *parser) parse() node {
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
			ast := parser.parse()
			result := ast.eval()

			fmt.Printf("%s = %d\n", eq, result)
		}(equation)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
