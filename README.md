# Parser for mathematical parser in Go
## What it does?
This Go application is a top-down parser for mathematical operations:
- `+`, `-`
- `*`, `/`
- `()`

The parser takes into account precedence and parenthesis. It evaluates equations by creating a tree, where a node can be a number or operation. Example tree for the equation: `(4 + 5 * (7 - 3)) - 2`,

## Assumptions:
- only one-digit numbers are allowed
- fractional numbers are not allowed
- machine stack is large enough.

To speed up execution, the program uses concurrency. 

As input program accepts the filename where every line contains an equation. If some equations are not valid program doesn't panic it returns an error message and proceeds with other lines (but it does panic in case of dividing by zero).

## How to run application

To run the application locally:

`go run parser.go <filename>`

To run tests:

`go test`

To run the application on docker:

`docker build -t math-parser .`

`docker run --rm math-parser`

It will use the example file `equations.txt` from this repo. If you want to pass your file on your local machine use the following:

`docker run --rm -v <path-to-your-local-file>:/app/equations.txt math-parser-go equations.txt`