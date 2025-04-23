package main

import (
	"fmt"
	"glox/ast"
	"glox/literal"
	"glox/token"
)

// Two jobs
// Produce a syntax tree from tokens
// Detect errors in a sequence of tokens

// A parser shoukd (with context of an errors)
// Detect and report the error
// If an error is not caught, bad stuff can happen in the back end
// Avoid crashing
// Valid input shoukd not cause it to loop infinently, input != valid code

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
// We syncronize back up with the unexplored par of the program that isnt directly affetced by the error
// This means we throw away all tokens on that line
type Parser struct {
	Tokens  []*token.Token
	Current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

func (p *Parser) parse() []ast.Stmt {
	statements := []ast.Stmt{}
	for !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			fmt.Println("Error in decl, syncronizing...")
			p.syncronize()
			continue
		}
		statements = append(statements, decl)
	}
	return statements
}

// Called repeatably to parse a series of statments in a program, perfect place to look for panic
func (p *Parser) declaration() (ast.Stmt, error) {
	// If current token is var, we are looking at var decl
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	if p.match(token.FUN) {
		return p.funDeclaration("function")
	}
	// If not, fallback to standard stmt
	return p.statement()
}

func (p *Parser) funDeclaration(kind string) (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name", kind))
	if err != nil {
		return nil, err
	}

	var params []token.Token
	if !p.match(token.RIGHT_PAREN) { // The next item is an identifier
		for {
			if len(params) >= 255 {
				ParserError(p.peek(), "Number of params exceeds 255 limit.")
			}
			val, err := p.consume(token.IDENTIFIER, "Expect idemtifier as paramteter.")
			if err != nil {
				return nil, err
			}
			params = append(params, val)
			if !p.match(token.COMMA) { break } // Break if we DONT see a comma
		}
	}

	// _, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	// if err != nil { 
	// 	return nil, err
	// }
	_, err = p.consume(token.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return ast.NewFun(name, params, body), nil
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	// Consume name only if the next token is an ident
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	// Set as nil empty expr
	var initializer ast.Expr
	if p.match(token.EQUAL) { // Look at the expression to give to a new var
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	// If = does not exist in decl, is empty decl, pass empty initial to var decl
	p.consume(token.SEMICOLON, "Expect ';' after variable in declaration.")
	return ast.NewVar(name, initializer), nil
}

// Decide what kind of statement to branch to
func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStmt()
	}

	if p.match(token.FOR) {
		return p.forStmt()
	}

	if p.match(token.IF) {
		return p.ifStmt()
	}

	if p.match(token.WHILE) {
		return p.whileStmt()
	}

	// Start of block statement
	if p.match(token.LEFT_BRACE) {
		block, err := p.block()
		if err != nil {
			return nil, err
		}

		return ast.NewBlock(block), nil
	}

	return p.exprStmt()
}

func (p *Parser) block() ([]ast.Stmt, error) {
	var stmts []ast.Stmt
	fmt.Println("Here")

	// While the next tok is not right brace and we are not at the end
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		fmt.Println(p.peek().Lexeme)
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, decl)
	}

	// The loop has concluded
	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return stmts, nil
}

func (p *Parser) forStmt() (ast.Stmt, error) {
	var err error
	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	// Dont forget
	// Match advances 'consumes' the next token if matched
	// Check returns wether the next is it or not simply

	var initializer ast.Stmt
	if p.match(token.SEMICOLON) { // Just a semicolon, this is directly following the opening (
		initializer = nil
	} else if p.match(token.VAR) { // Init is a new var, x := 1
		if initializer, err = p.varDeclaration(); err != nil {
			return nil, err
		}
	} else { // Is expression, hopefully with side effect
		if initializer, err = p.exprStmt(); err != nil {
			return nil, err
		}
	}

	var condition ast.Expr = nil
	if !p.check(token.SEMICOLON) { // See if next token is not semicolon, dont consume it
		if condition, err = p.expression(); err != nil {
			return nil, err // If not, parse expression, not matter what ';' should be at end
		}
	}
	// Consume now
	_, err = p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	// Same here but we expect a closing paren instead
	var increment ast.Expr = nil
	if !p.check(token.RIGHT_PAREN) {
		if increment, err = p.expression(); err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after 'for' clauses.")
	if err != nil {
		return nil, err
	}

	var body ast.Stmt = nil
	if body, err = p.statement(); err != nil {
		return nil, err
	}
	if increment != nil {
		body = ast.NewBlock([]ast.Stmt{body, ast.NewExpression(increment)})
	}
	if condition == nil {
		fmt.Println("Not run")
		condition = ast.NewLiteralExpr(literal.NewLiteral(true))
	}
	body = ast.NewWhile(condition, body)
	if initializer != nil {
		body = ast.NewBlock([]ast.Stmt{initializer, body})
	}

	return body, nil
}

func (p *Parser) whileStmt() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_BRACE, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return ast.NewWhile(condition, body), nil
}

func (p *Parser) ifStmt() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_BRACE, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var finalBranch ast.Stmt

	if p.match(token.ELSE) {
		finalBranch, err = p.statement()
	}

	return ast.NewIf(condition, thenBranch, finalBranch), nil
}

func (p *Parser) printStmt() (ast.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return ast.NewPrint(val), nil
}

func (p *Parser) exprStmt() (ast.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return ast.NewExpression(val), nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	var val ast.Expr

	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		// equals := p.previous()
		if val, err = p.assignment(); err != nil {
			return nil, err
		}

		// If expr is a var
		if varExpr, ok := expr.(*ast.VarExpr); ok {
			name := varExpr.Name
			return ast.NewAssignExpr(name, val), nil
		}
	}


	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = ast.NewLogicalExpr(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
	var err error
	var expr ast.Expr

	if expr, err = p.equality(); err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality() // Calling equality then begins seeing the boolean val of the left expression
		if err != nil {
			return nil, err
		}
		expr = ast.NewLogicalExpr(expr, operator, right)
	}
	return expr, nil
}

// If we are parsing a == b == c
// We parse a == b, then a == b becomes the left operand of == c, looping
// Returning the entire expression at the end
func (p *Parser) equality() (ast.Expr, error) {
	expr, _ := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		fmt.Println("Equality prev lexeme: ", operator.Lexeme)

		right, _ := p.comparison()
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, _ := p.term()

	// While we are currently in a token that is composed of 2 of these
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, _ := p.term()
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, _ := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, _ := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.BANG, token.MINUS) { // If it is ! or -, must be unary
		operator := p.previous()
		right, err := p.unary() // Parse recursively, ie, !!
		if err != nil {
			return nil, err
		}
		return ast.NewUnaryExpr(operator, right), nil
	}

	return p.call()
	// Must have reached highest level precedence
}

func (p *Parser) call() (ast.Expr, error) {
	var expr ast.Expr
	var err error
	if expr, err = p.primary(); err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) { // If, after consuming maybe identifier, an opening paren exists
			if expr, err = p.finishCall(expr); err != nil {
				return nil, err
			} else {
				break
			}
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	var args []ast.Expr

	if !p.check(token.RIGHT_PAREN) { // If we dont see right paren as we are walking the args
		// PArse expr
		// Then see if there is a comma
		// No right paren? Expect expr (Arg)
		for {
			if len(args) >= 255 {
				fmt.Println("Args is over 255!")
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)

			if !p.match(token.COMMA) {
				break
			}
			_, err = p.consume(token.COMMA, "Expect ',' in argument list.")
			if err != nil {
				return nil, err
			}
			// Should we consume the comma if its there?
		}
	}

	paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	// Finally perform func call
	return ast.NewCallExpr(callee, paren, args), nil
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return ast.NewLiteralExpr(literal.NewLiteral(false)), nil
	}

	if p.match(token.TRUE) {
		return ast.NewLiteralExpr(literal.NewLiteral(true)), nil
	}

	if p.match(token.NIL) {
		return ast.NewLiteralExpr(literal.NewLiteral(nil)), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		fmt.Println("str")
		return ast.NewLiteralExpr(p.previous().Literal), nil
	}

	// If we see an ident
	if p.match(token.IDENTIFIER) {
		return ast.NewVarExpr(p.previous()), nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return ast.NewGroupingExpr(expr), nil
	}

	// If we have gotten here, the token given cannot start an expression
	ParserError(p.peek(), "Expect expression.")
	return nil, nil
}

// If passed token is the type of next token, consume it, otherwise error
func (p *Parser) consume(tokType token.TokenType, message string) (token.Token, error) {
	if p.check(tokType) {
		return p.advance(), nil
	} // If next token is passed type, consume it and pass the previous token

	ParserError(p.peek(), message)
	return token.Token{}, fmt.Errorf(message)
}

// Discards tokens until it has found the end of a statement
// Now we begin again at the next statement
// Hopefully all tokens that would have been affected by an earlier error are discorded
func (p *Parser) syncronize() {
	p.advance() // Consume a token

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		} // Found statement boundary

		switch p.peek().Type {
		case token.CLASS:
		case token.FUN:
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
	return *p.Tokens[p.Current]
}

func (p *Parser) previous() token.Token {
	return *p.Tokens[p.Current-1]
}

func (p *Parser) advance() token.Token { // Returns token that is consumed, s.Current-1
	if !p.isAtEnd() {
		p.Current += 1
	} // Consume as long as we are not at the end
	return p.previous()
}

// Got or to work
// Need to get level or precedence explained
