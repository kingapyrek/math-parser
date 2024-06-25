package main

import (
	"fmt"
	"strconv"
)

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
	expression := "(4 + 5 * (7 - 3)) - 2"
	parser := &parser{expression: expression, pos: 0}

	//start parsing expression
	ast := parser.parse()

	// Evaluation is a separate step
	result := ast.eval()
	fmt.Println(result) // Output: 22
}
