package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	HadError    bool
	Interpreter *Interpreter
}

func main() {
	// expression := expression.NewBinaryExpr(
	// 	expression.NewLiteralExpr(literal.NewFloatLiteral(123)), *token.NewToken(token.MINUS, "-", nil, 1), expression.NewLiteralExpr(literal.NewFloatLiteral(123)))

	// fmt.Println(ast.NewAstPrinter().Print(expression))

	lox := NewLox()
	lox.Start()
}

func NewLox() *Lox {
	return &Lox{
		HadError:    false,
		Interpreter: NewInterpreter(),
	}
}

func (l *Lox) Start() {
	args := os.Args

	if len(args) > 2 {
		fmt.Println("Usage glox [script]")
		os.Exit(64) // Incorrect usage
	} else if len(args) == 2 {
		fmt.Println("Running file")
		l.RunFile(args[1])
	} else {
		l.RunPrompt()
	}
}

func (l *Lox) RunFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		os.Exit(74)
	}
	l.Run(string(data))

	if l.HadError {
		os.Exit(65)
	}
}

func (l *Lox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		l.Run(line)
		l.HadError = false // Reset so it does not kill the session
	}
}

func (l *Lox) Run(source string) {
	// Parse characters into tokens
	scanner := NewGloxScanner(source, l)
	tokens := scanner.ScanTokens()

	// Parse tokens to expressions
	parser := NewParser(tokens)
	statements := parser.parse()

	if l.HadError {
		fmt.Println("Error encountered in Run")
		return
	}

	if l.Interpreter.HadRuntimeError {
		fmt.Println("Runtime Error encountered in Run")
		return
	}

	// fmt.Println(ast.NewAstPrinter().Print(expr))

	l.Interpreter.interpret(statements)
}
