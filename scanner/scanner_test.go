package scanner

import (
	"glorp/token"
	"testing"
)

func TestPeek(t *testing.T) {
	source := "var x = 42"
	scanner := NewScanner()
	scanner.Source = source

	// Test initial peek
	ch := scanner.peek()
	if ch != 'v' {
		t.Errorf("Expected 'v', got '%c'", ch)
	}

	// Advance and peek again
	scanner.advance()
	ch = scanner.peek()
	if ch != 'a' {
		t.Errorf("Expected 'a', got '%c'", ch)
	}

	for !scanner.isAtEnd() {
		scanner.advance()
	}

	scanner.advance()
	ch = scanner.peek()
	if ch != 0 {
		t.Errorf("Expected 0, got '%c'", ch)
	}
}

func TestNumber(t *testing.T) {
	source1 := "100"
	source2 := "20.4"
	scanner := NewScanner()
	
	scanner.Source = source1
	scanner.number()
	if scanner.Tokens[0].Type != token.NUMBER {
		t.Errorf("Expected token to be lex of type Number but was not")
	} 
	if scanner.Tokens[0].Lexeme != "100" {
		t.Errorf("Expected token to be lex of 100, got %s", scanner.Tokens[0].Lexeme)
	}

	scanner.Source = source2
	scanner.number()
	if scanner.Tokens[1].Type != token.NUMBER {
		t.Errorf("Expected token to be lex of type Number but was not")
	}
	if scanner.Tokens[1].Lexeme != "20.4" {
		t.Errorf("Expected token to be lex of 20.4, got %s", scanner.Tokens[1].Lexeme)
	}
}

func TestFutureChar(t *testing.T) {
	source := "     }"
	scanner := NewScanner()
	scanner.Source = source

	if scanner.futureChar() != '}' {
		t.Errorf("Expected future of '}', got %v", scanner.futureChar())
	}
}