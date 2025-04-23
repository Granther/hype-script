package token

import (
	"fmt"
	"hype-script/internal/literal"
)

type TokenType int

const (
	// Single-character tokens.
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_BRACKET
	RIGHT_BRACKET
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	END
	SPACE

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	PLUS_EQUAL
	PLUS_PLUS
	MINUS_MINUS
	MINUS_EQUAL
	STAR_EQUAL
	SLASH_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	GLUNC
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	WERT
	WOOPS
	TRY
	IMPORT
	AS

	// End of file
	EOF
)

var TokenTypeNames = map[TokenType]string{
	LEFT_PAREN:    "LEFT_PAREN",
	RIGHT_PAREN:   "RIGHT_PAREN",
	LEFT_BRACE:    "LEFT_BRACE",
	RIGHT_BRACE:   "RIGHT_BRACE",
	COMMA:         "COMMA",
	DOT:           "DOT",
	MINUS:         "MINUS",
	PLUS:          "PLUS",
	SEMICOLON:     "SEMICOLON",
	SLASH:         "SLASH",
	STAR:          "STAR",
	END:           "END",
	SPACE:         "SPACE",
	BANG:          "BANG",
	BANG_EQUAL:    "BANG_EQUAL",
	EQUAL:         "EQUAL",
	EQUAL_EQUAL:   "EQUAL_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",
	IDENTIFIER:    "IDENTIFIER",
	STRING:        "STRING",
	NUMBER:        "NUMBER",
	AND:           "AND",
	CLASS:         "CLASS",
	ELSE:          "ELSE",
	TRUE:          "TRUE",
	FALSE:         "FALSE",
	GLUNC:         "GLUNC",
	FOR:           "FOR",
	IF:            "IF",
	NIL:           "NIL",
	OR:            "OR",
	PRINT:         "PRINT",
	RETURN:        "RETURN",
	WERT:          "WERT",
	WOOPS:         "WOOPS",
	TRY:           "TRY",
	VAR:           "VAR",
	PLUS_EQUAL:    "PLUS_EQUAL",
	PLUS_PLUS:     "PLUS_PLUS",
	MINUS_MINUS:   "MINUS_MINUS",
	MINUS_EQUAL:   "MINUS_EQUAL",
	STAR_EQUAL:    "STAR_MINUS",
	SLASH_EQUAL:   "SLASH_EQUAL",
	LEFT_BRACKET:  "LEFT_BRACKET",
	RIGHT_BRACKET: "RIGHT_BRACKET",
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal *literal.Literal
	Line    int
}

func NewToken(tokType TokenType, lexeme string, literal *literal.Literal, line int) *Token {
	return &Token{
		Type:    tokType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%d %s %s", t.Type, t.Lexeme, t.Literal.String())
}

func BuildKeywords() (keywords map[string]TokenType) {
	keywords = make(map[string]TokenType)
	keywords["import"] = IMPORT
	keywords["as"] = AS
	keywords["and"] = AND
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["for"] = FOR
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["true"] = TRUE
	keywords["while"] = WHILE
	keywords["try"] = TRY
	return
}

func BuildLeftOper() (leftOperators map[rune]TokenType) {
	leftOperators = make(map[rune]TokenType)
	leftOperators['+'] = PLUS
	leftOperators['-'] = MINUS
	leftOperators['*'] = STAR
	leftOperators['/'] = SLASH
	leftOperators['='] = EQUAL
	return
}
