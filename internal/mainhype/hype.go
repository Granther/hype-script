package mainhype

import (
	"os"
	"path/filepath"
	"bufio"
	"fmt"
	"hype-script/internal/environment"
	"hype-script/internal/interpreter"
	"hype-script/internal/parser"
	"hype-script/internal/scanner"
	"hype-script/internal/token"
	"hype-script/internal/types"
)

type Hype struct {
	HadError    bool
	Scanner     *scanner.Scanner
	Parser      types.Parser
	Interpreter types.Interpreter
	Environment types.Environment
}

func NewHype() *Hype {
	env := environment.NewEnvironment(nil)
	return &Hype{
		HadError:    false,
		Scanner:     scanner.NewScanner(),
		Parser:      parser.NewParser(env),
		Interpreter: interpreter.NewInterpreter(env),
		Environment: env,
	}
}

func (g *Hype) Start() error {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: hype [file.hyp]")
		return nil
	} else if len(args) == 2 {
		return g.Runfile(args[1])
	} else {
		return g.Repl()
	}
}

func (g *Hype) Runfile(file string) error {
	ext := filepath.Ext(file)
	if ext != ".hyp" {
		return fmt.Errorf("hype file (.hyp) is required to run, got %s", ext)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return g.Run(string(data))
}

func (g *Hype) Repl() error {
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

func (g *Hype) Run(source string) error {
	tokens, err := g.Scanner.ScanTokens(source)
	if err != nil {
		return err
	}

	// Debug to see token types
	for _, tok := range tokens {
		fmt.Printf("%s %s\n", token.TokenTypeNames[tok.Type], tok.Lexeme)
	}

	statements := g.Parser.Parse(tokens)
	// Debug to see statement info
	// for _, stmt := range statements {
	// 	fmt.Printf("%s\n", stmt.String())
	// }
	//fmt.Printf("Num stmts: %d\n", len(statements))
	
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