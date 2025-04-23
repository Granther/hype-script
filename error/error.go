package error

import (
	"fmt"
	"glorp/token"
)

func ParserError(errToken token.Token, message string) {
	if errToken.Type == token.EOF {
		Report(errToken.Line, " at end", message, "parser")
	} else {
		Report(errToken.Line, fmt.Sprintf("at '%s'", errToken.Lexeme), message, "parser")
	}
}

func InterpreterRuntimeError(errToken token.Token, message string) {
	Report(errToken.Line, fmt.Sprintf(" at '%s'", errToken.Lexeme), message, "interpreter")
}

func InterpreterSimpleRuntimeError(errToken token.Token, message string) {
	Report(errToken.Line, fmt.Sprintf(" at '%s'", errToken.Lexeme), message, "interpreter")
}

func ScannerError(line int, message string) {
	Report(line, "", message, "scanner")
}

func Report(line int, where string, message string, subsystem string) {
	fmt.Printf("[subsystem %s] [line %d] Error %s: %s\n", subsystem, line, where, message)
}
