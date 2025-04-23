package scanner

import (
	"fmt"
	glorpError "glorp/error"
	"glorp/literal"
	"glorp/token"
	"strconv"
)

type Scanner struct {
	Source        string
	Tokens        []token.Token
	Start         int // Points to first lexeme being scanner
	Current       int // Character currently being considered
	Line          int // The source line that Current is on
	Keywords      map[string]token.TokenType
	LeftOperators map[rune]token.TokenType
}

func NewScanner() *Scanner {
	keywords := make(map[string]token.TokenType)
	keywords["and"] = token.AND
	keywords["class"] = token.CLASS
	keywords["else"] = token.ELSE
	keywords["false"] = token.FALSE
	keywords["for"] = token.FOR
	keywords["glunc"] = token.GLUNC
	keywords["if"] = token.IF
	keywords["nil"] = token.NIL
	keywords["or"] = token.OR
	keywords["print"] = token.PRINT
	keywords["return"] = token.RETURN
	keywords["super"] = token.SUPER
	keywords["this"] = token.THIS
	keywords["true"] = token.TRUE
	keywords["var"] = token.VAR
	keywords["while"] = token.WHILE
	keywords["wert"] = token.WERT
	keywords["try"] = token.TRY
	keywords["woops"] = token.WOOPS

	leftOperators := make(map[rune]token.TokenType)
	leftOperators['+'] = token.PLUS
	leftOperators['-'] = token.MINUS
	leftOperators['*'] = token.STAR
	leftOperators['/'] = token.SLASH
	leftOperators['='] = token.EQUAL

	return &Scanner{
		Tokens:        []token.Token{},
		Start:         0,
		Current:       0,
		Line:          1,
		Keywords:      keywords,
		LeftOperators: leftOperators,
	}
}

func (s *Scanner) ScanTokens(source string) ([]token.Token, error) {
	s.Source = source
	// Each iteration we scan a single token
	for !s.isAtEnd() {
		s.Start = s.Current
		s.scanToken()
	}

	if s.Tokens[len(s.Tokens)-1].Type != token.END {
		endToken := token.NewToken(token.END, "", nil, s.Line)
		s.Tokens = append(s.Tokens, *endToken)
	}

	// Appends an EOF token at the end
	newToken := token.NewToken(token.EOF, "", nil, s.Line)
	s.Tokens = append(s.Tokens, *newToken)

	return s.Tokens, nil
}

// If current character being checked is >= len of source
// If we are parsing beyond the source return true
func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	c := s.advance()

	// fmt.Println(string(c))

	switch c {
	case '(':
		s.addSimpleToken(token.LEFT_PAREN)
	case ')':
		s.addSimpleToken(token.RIGHT_PAREN)
		if s.futureChar() == '}' {
			s.addSimpleToken(token.END)
		}
	case '{':
		s.addSimpleToken(token.LEFT_BRACE)
	case '}':
		s.addSimpleToken(token.RIGHT_BRACE)
		for s.peek() == ' ' {
			s.advance()
		}
		if s.peek() != '\n' {
			s.addSimpleToken(token.END)
		}
	case ',':
		s.addSimpleToken(token.COMMA)
	case '.':
		s.addSimpleToken(token.DOT)
	case '-':
		if s.match('=') {
			s.addSimpleToken(token.MINUS_EQUAL)
		} else if s.match('-') {
			s.addSimpleToken(token.MINUS_MINUS)
		} else {
			s.addSimpleToken(token.MINUS)
		}
	case '+':
		if s.match('=') {
			s.addSimpleToken(token.PLUS_EQUAL)
		} else if s.match('+') {
			s.addSimpleToken(token.PLUS_PLUS)
		} else {
			s.addSimpleToken(token.PLUS)
		}
	case ';':
		s.addSimpleToken(token.END)
	case '*':
		if s.match('=') {
			s.addSimpleToken(token.STAR_EQUAL)
		} else {
			s.addSimpleToken(token.STAR)
		}
	case '=':
		if s.match('=') {
			s.addSimpleToken(token.EQUAL_EQUAL)
		} else {
			s.addSimpleToken(token.EQUAL)
		}
	case '>':
		if s.match('=') {
			s.addSimpleToken(token.GREATER_EQUAL)
		} else {
			s.addSimpleToken(token.GREATER)
		}
	case '<':
		if s.match('=') {
			s.addSimpleToken(token.LESS_EQUAL)
		} else {
			s.addSimpleToken(token.LESS)
		}
	case '!': // Are we looking at a lexeme of ! OR !=
		if s.match('=') {
			s.addSimpleToken(token.BANG_EQUAL)
		} else {
			s.addSimpleToken(token.BANG)
		}
	case '/': // Are we doing division or commenting?
		if s.match('/') { // If next char is /, is comment, read till the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				// for !s.match('\n') && !s.isAtEnd() {
				s.advance()
			}
			if s.match('\n') {
				s.Line++
			}
		} else if s.match('=') {
			s.addSimpleToken(token.SLASH_EQUAL)
		} else {
			s.addSimpleToken(token.SLASH)
		}
	case ' ': // We are basically skipping these, no error, no op
	case '\r':
	case '\t':
		// Ignore whitespace
	case '\n': // Do nothing but start iterate to the next line
		s.addSimpleToken(token.END)
		s.Line += 1
		for s.match('\n') {
			s.Line += 1
		}
	case '[':
		// s.glist()
		s.addSimpleToken(token.LEFT_BRACKET)
	case ']':
		s.addSimpleToken(token.RIGHT_BRACKET)
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			glorpError.ScannerError(s.Line, "Unexpected character")
		}
	}
}

// Similar to advance but does not consume the character, 'lookahead'
func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return rune(s.Source[s.Current])
}

func (s *Scanner) nextIsOper() bool {
	if s.peek() == ' ' {
		s.advance()
	}
	_, ok := s.LeftOperators[s.peek()]
	return ok
}

func (s *Scanner) attemptEarlyEnd() {
	if !s.nextIsOper() {
		s.addSimpleToken(token.END)
		s.match('\n')
	}
}

func (s *Scanner) futureChar() rune {
	cur := s.Current
	for s.peek() == ' ' {
		s.Current++
	}
	p := s.peek()
	s.Current = cur
	return p
}

func (s *Scanner) prev() rune {
	return rune(s.Source[s.Current-1])
}

func (s *Scanner) peekNext() rune {
	// If current + 1 is greater if equal to len of source, if source is 10, and current is 10
	if s.Current+1 >= len(s.Source) {
		return '0'
	}
	return rune(s.Source[s.Current+1])
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() { // Keep searching for string closing
		if s.peek() == '\n' {
			s.Line += 1
		}
		s.advance()
	}

	s.advance()

	if s.isAtEnd() { // If it makes it to the end of line before finding closing "
		glorpError.ScannerError(s.Line, "Unterminated string")
		return
	}

	// s.Current is one less than the closing ", make Current into "
	// s.advance()

	// Cut the begining and closing "'s off
	val := s.Source[s.Start+1 : s.Current-1]
	lit := literal.NewLiteral(val)
	s.addToken(token.STRING, lit)

	// s.attemptEarlyEnd()
}

// Consumes next character of source line and returns it
func (s *Scanner) advance() rune {
	if s.isAtEnd() {
		return '0'
	}
	// I give you Grant, the dumbest person alive...
	sub := s.Source[s.Current]
	s.Current += 1
	return rune(sub)
}

func (s *Scanner) addSimpleToken(tokType token.TokenType) {
	s.addToken(tokType, nil)
}

func (s *Scanner) addToken(tokType token.TokenType, literal *literal.Literal) {
	text := s.Source[s.Start:s.Current]
	escapedText := strconv.QuoteToASCII(text)
	escapedText = escapedText[1 : len(escapedText)-1] // Remove the surrounding quotes added by QuoteToASCII
	newToken := token.NewToken(tokType, escapedText, literal, s.Line)
	s.Tokens = append(s.Tokens, *newToken)
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	} // There is no next token, we are at the end
	if s.Source[s.Current] != expected {
		return false
	} // If it is not expected, ie, = after !, return false

	s.Current += 1 // We are done looking at this char only if the next char is what we expected
	return true    // The passed char is what we expected
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	// While the characters being explored are part d
	for s.isDigit(s.peek()) {
		s.advance() // What if we try to advance but are at the end?
	}

	// If number ends in ., dont parse .
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the '.'
		s.advance()

		// Then parse the rest, after the dot
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	f64, _ := strconv.ParseFloat(s.Source[s.Start:s.Current], 32)
	f64Literal := literal.NewLiteral(f64)
	s.addToken(token.NUMBER, f64Literal)

	// s.attemptEarlyEnd()
}

func (s *Scanner) identifier() {
	// While next char is alphanumeric, advance
	for s.isAlphaNumeric(s.peek()) && !s.isAtEnd() {
		s.advance()
	}

	text := s.Source[s.Start:s.Current]

	tokType, ok := s.Keywords[text]
	if !ok { // If it is not a recognized keyword, label it as ident
		s.addSimpleToken(token.IDENTIFIER)
		return
	}

	s.addSimpleToken(tokType)

	if s.futureChar() == '}' {
		fmt.Println("future char")
		s.addSimpleToken(token.END)
	}
}

// Check and see if a byte char is alpha numeric
func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// If is alphabetical or number
func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}
