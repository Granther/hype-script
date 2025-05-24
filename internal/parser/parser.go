package parser

import (
	herror "hype-script/internal/error"
	"hype-script/internal/token"
	"hype-script/internal/types"
	"errors"
)

// Two jobs
// Produce a syntax tree from tokens
// Detect errors in a sequence of tokens

// A parser should (with context of an errors)
// Detect and report the error
// If an error is not caught, bad stuff can happen in the back end
// Avoid crashing
// Valid input should not cause it to loop infinently, input != valid code

// Minimize cascaded errors
// Errors that occur as a result of earlier errors
// These errors are only results of the initial error

// Panic mode recovery
// Able to jump out of an expression that contains an error to reduce cascaded errors

// Adding an error to the syntax of the language allows of graceful handling of common errors
// This also allows for better, more specific error messages

// Parser's state
// What rules it is in the middle of parsing

// Syncronizing after a panic
// We syncronize back up with the unexplored part of the program that isnt directly affected by the error
// This means we throw away all tokens on that line

type Parser struct {
	HadError    bool
	Tokens      []token.Token
	Environment types.Environment
	Current     int
}

func NewParser(e types.Environment) *Parser {
	return &Parser{
		HadError:    false,
		Current:     0,
		Environment: e,
	}
}

// Takes in parsed tokens from Scanner and outputs list of Statements
func (p *Parser) Parse(tokens []token.Token) []types.Stmt {
	p.Tokens = tokens
	statements := []types.Stmt{}

	// We see no tokens, just return
	if len(tokens) == 0 { return statements }

	// While we are still within range of passed tokens
	for !p.isAtEnd() {
		p.match(token.END)           // Consume endline token if its there
		decl, err := p.declaration() // Decl is start of recursive statment parsing
		if err != nil {
			p.HadError = true
			p.syncronize()
			continue
		}
		statements = append(statements, decl)
	}
	return statements
}

// If passed token is the type of next token, consume it, otherwise error
func (p *Parser) consume(tokType token.TokenType, message string) (token.Token, error) {
	if p.check(tokType) {
		return p.advance(), nil
	} // If next token is passed type, consume it and pass the previous token

	herror.ParserError(p.peek(), message)
	return token.Token{}, errors.New(message)
}

// Discards tokens until it has found the end of a statement
// Now we begin again at the next statement
// Hopefully all tokens that would have been affected by an earlier error are discorded
func (p *Parser) syncronize() {
	p.advance() // Consume a token

	for !p.isAtEnd() {
		if p.previous().Type == token.END {
			return
		} // Found statement boundary

		switch p.peek().Type {
		case token.FUN:
		case token.PAR:
		case token.HYP:
		case token.VAR:
		case token.FOR:
		case token.IF:
		case token.WHILE:
		case token.PRINT:
		case token.RETURN: // Found statement boundry here too
			return
		}
		p.advance()
	}
}

// Recovery is a big deal in the parser

// See if the passed tokens, if we dont see one of them, we must be done
func (p *Parser) match(tokenTypes ...token.TokenType) bool {
	for _, tokType := range tokenTypes {
		if p.check(tokType) { // Only consumes token if it/the next is what it is looking for
			p.advance()
			return true
		}
	}
	return false
}

// Only looks at tokens
// Check if the next token is the passed type
func (p *Parser) check(tokType token.TokenType) bool {
	if p.isAtEnd() { // Because there is no next token
		return false
	}
	return p.peek().Type == tokType // Does the next token == the passed one?
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token { // Returns current token we have yet to consume
	return p.Tokens[p.Current]
}

func (p *Parser) peekNext() token.Token { // Returns current token we have yet to consume
	if p.isAtEnd() {
		return token.Token{}
	}
	return p.Tokens[p.Current+1]
}

func (p *Parser) previous() token.Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) advance() token.Token { // Returns token that is consumed, s.Current-1
	if !p.isAtEnd() {
		p.Current += 1
	} // Consume as long as we are not at the end
	return p.previous()
}

func (p *Parser) GetHadError() bool {
	return p.HadError
}

// func (p *Parser) matchNext(tokenTypes ...token.TokenType) bool {
// 	if p.isAtEnd() {
// 		return false
// 	}
// 	for _, tokType := range tokenTypes {
// 		if p.checkNext(tokType) {
// 			p.advance()
// 			return true
// 		}
// 	}
// 	return false
// }

// func (p *Parser) checkNext(tokType token.TokenType) bool {
// 	if p.isAtEnd() { // Because there is no next token
// 		return false
// 	}
// 	return p.peekNext().Type == tokType
// }

// func (p *Parser) classDeclaration() (types.Stmt, error) {
// 	name, err := p.consume(token.IDENTIFIER, "Expect a valid name following 'class'.")
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = p.consume(token.LEFT_BRACE, "Expect '{' after class name.")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var methods []types.Stmt
// 	for !p.match(token.RIGHT_BRACE) {
// 		if p.match(token.END) { continue }
// 		method, err := p.funDeclaration("method")
// 		if err != nil {
// 			return nil, err
// 		}
// 		methods = append(methods, method)
// 	}

// 	fmt.Print(methods)

// 	return types.NewClass(name, methods), nil
// }

// A class is a list of func decls with a name, a named list of func decls!
