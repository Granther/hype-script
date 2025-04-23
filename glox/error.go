package main

import (
	"fmt"
	"glox/token"
)

func ParserError(errToken token.Token, message string) {
	if errToken.Type == token.EOF {
		Report(errToken.Line, " at end", message)
	} else {
		Report(errToken.Line, fmt.Sprintf(" at '%s'", errToken.Lexeme), message)
	}
}

func InterpreterRuntimeError(errToken token.Token, message string) {
	Report(errToken.Line, fmt.Sprintf(" at '%s'", errToken.Lexeme), message)
}

func (l *Lox) Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
}
