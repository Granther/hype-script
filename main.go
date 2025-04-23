package main

import (
	"bufio"
	"fmt"
	"glorp/environment"
	"glorp/interpreter"
	"glorp/parser"
	"glorp/scanner"
	"glorp/token"
	"glorp/types"
	"os"
	"path/filepath"
)

type Glorp struct {
	HadError    bool
	Scanner     *scanner.Scanner
	Parser      types.Parser
	Interpreter types.Interpreter
	Environment types.Environment
}

func NewGlorp() *Glorp {
	env := environment.NewEnvironment(nil)
	return &Glorp{
		HadError:    false,
		Scanner:     scanner.NewScanner(),
		Parser:      parser.NewParser(env),
		Interpreter: interpreter.NewInterpreter(env),
		Environment: env,
	}
}

func (g *Glorp) Start() error {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: glorp [file.glp]")
		return nil
	} else if len(args) == 2 {
		return g.Runfile(args[1])
	} else {
		return g.Repl()
	}
}

func (g *Glorp) Runfile(file string) error {
	ext := filepath.Ext(file)
	if ext != ".glp" {
		return fmt.Errorf("glorp file (.glp) is required to run, got %s", ext)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return g.Run(string(data))
}

func (g *Glorp) Repl() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		g.Run(line)
		g.HadError = false
	}
	return nil
}

func (g *Glorp) Run(source string) error {
	tokens, err := g.Scanner.ScanTokens(source)
	if err != nil {
		return err
	}

	for _, tok := range tokens {
		fmt.Printf("%s %s\n", token.TokenTypeNames[tok.Type], tok.Lexeme)
	}

	statements := g.Parser.Parse(tokens)

	if g.Parser.GetHadError() {
		fmt.Println("Error encountered in Parser, stopping...")
		return nil
	}

	g.Interpreter.Interpret(statements)

	if g.Interpreter.GetHadRuntimeError() {
		fmt.Println("Runtime Error encountered in Run")
		return nil
	}

	return nil
}

func main() {
	glorp := NewGlorp()
	err := glorp.Start()
	if err != nil {
		fmt.Println("Unable to GLORP: ", err)
		os.Exit(1)
	}
}
