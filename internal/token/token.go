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
	TILDE // ~
	KARAT // ^

	// One or two character tokens.
	BANG       // !
	BANG_EQUAL // !=
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
	TILDE_EQUAL // ~=

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	ELSE
	FALSE
	FUN
	FOR
	IF
	NEWT // None
	OR
	PRINT // Print stmt before Std lib
	RETURN
	TRUE
	VAR
	WHILE
	IMPORT
	AS
	PAR // Parallel
	HYP // Hype

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
	TILDE:         "TILDE",
	TILDE_EQUAL:   "TILDE_EQUAL",
	KARAT:         "KARAT",
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
	ELSE:          "ELSE",
	TRUE:          "TRUE",
	FALSE:         "FALSE",
	FUN:           "FUN",
	FOR:           "FOR",
	IF:            "IF",
	NEWT:          "NEWT",
	OR:            "OR",
	PRINT:         "PRINT",
	RETURN:        "RETURN",
	VAR:           "VAR",
	PAR:           "PAR",
	HYP:           "HYP",
	PLUS_EQUAL:    "PLUS_EQUAL",
	PLUS_PLUS:     "PLUS_PLUS",
	MINUS_MINUS:   "MINUS_MINUS",
	MINUS_EQUAL:   "MINUS_EQUAL",
	STAR_EQUAL:    "STAR_MINUS",
	SLASH_EQUAL:   "SLASH_EQUAL",
	LEFT_BRACKET:  "LEFT_BRACKET",
	RIGHT_BRACKET: "RIGHT_BRACKET",
	IMPORT:        "IMPORT",
}

var BadTokens = map[rune]bool{
	' ':  true,
	'\t': true,
	'\r': true,
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

// Gives the token type, string value used to identify it, and literal val
func (t *Token) String() string {
	lit := "nil"
	if t.Literal != nil {
		lit = t.Literal.String()
	}
	return fmt.Sprintf("Type: %s Lex: %s Literal: %s", TokenTypeNames[t.Type], t.Lexeme, lit)
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
	keywords["newt"] = NEWT
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["true"] = TRUE
	keywords["while"] = WHILE
	keywords["func"] = FUN
	keywords["var"] = VAR
	keywords["par"] = PAR
	keywords["hyp"] = HYP
	return
}

func BuildLeftOper() (leftOperators map[rune]TokenType) {
	leftOperators = make(map[rune]TokenType)
	leftOperators['+'] = PLUS
	leftOperators['-'] = MINUS
	leftOperators['*'] = STAR
	leftOperators['/'] = SLASH
	leftOperators['='] = EQUAL
	leftOperators['~'] = TILDE
	leftOperators['^'] = KARAT
	return
}
