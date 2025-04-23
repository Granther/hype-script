package parser

import (
	"fmt"
	herror "hype-script/internal/error"
	"hype-script/internal/literal"
	"hype-script/internal/token"
	"hype-script/internal/types"
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

func (p *Parser) Parse(tokens []token.Token) []types.Stmt {
	p.Tokens = tokens

	statements := []types.Stmt{}
	for !p.isAtEnd() {
		p.match(token.END)
		decl, err := p.declaration()
		if err != nil {
			p.HadError = true
			p.syncronize()
			continue
		}
		statements = append(statements, decl)
	}
	return statements
}

// Called repeatably to parse a series of statments in a program, perfect place to look for panic
func (p *Parser) declaration() (types.Stmt, error) {
	// If current token is var, we are looking at var decl
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	if p.match(token.GLUNC) {
		return p.funDeclaration("function")
	}
	if p.match(token.CLASS) {
		return p.classDeclaration()
	}
	// If not, fallback to standard stmt
	return p.statement()
}

func (p *Parser) classDeclaration() (types.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect a valid name following 'class'.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.LEFT_BRACE, "Expect '{' after class name.")
	if err != nil {
		return nil, err
	}

	var methods []types.Stmt
	for !p.match(token.RIGHT_BRACE) {
		if p.match(token.END) { continue }
		method, err := p.funDeclaration("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	fmt.Print(methods)

	return types.NewClass(name, methods), nil
}

// A class is a list of func decls with a name, a named list of func decls!

func (p *Parser) funDeclaration(kind string) (types.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name", kind))
	if err != nil {
		return nil, err
	}

	var params []token.Token
	if !p.check(token.RIGHT_PAREN) { // The next item is an identifier
		for {
			if len(params) >= 255 {
				herror.ParserError(p.peek(), "Number of params exceeds 255 limit.")
			}
			val, err := p.consume(token.IDENTIFIER, "Expect identifier as paramteter.")
			if err != nil {
				return nil, err
			}
			params = append(params, val)

			if p.check(token.RIGHT_PAREN) {
				break
			}
			if !p.match(token.COMMA) {
				break
			} // Break if we DONT see a comma
		}
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return types.NewFun(name, params, body, p.Environment), nil
}

// The problem is that all functions are defined within the global scope
// So we must define each func within the

func (p *Parser) varDeclaration() (types.Stmt, error) {
	// Consume name only if the next token is an ident
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	// Set as nil empty expr
	var initializer types.Expr
	if p.match(token.EQUAL) { // Look at the expression to give to a new var
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.END, "Expect 'end' after var decl/init.")
	if err != nil {
		return nil, err
	}

	return types.NewVar(name, initializer), nil
	// If = does not exist in decl, is empty decl, pass empty initial to var decl
}

// Decide what kind of statement to branch to
func (p *Parser) statement() (types.Stmt, error) {
	if p.match(token.RETURN) {
		return p.returnStmt()
	}

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

	if p.match(token.WERT) {
		return p.wertStmt()
	}

	if p.match(token.TRY) {
		return p.tryStmt()
	}

	// Start of block statement
	if p.match(token.LEFT_BRACE) {
		block, err := p.block()
		if err != nil {
			return nil, err
		}

		return types.NewBlock(block), nil
	}

	return p.exprStmt()
}

func (p *Parser) block() ([]types.Stmt, error) {
	var stmts []types.Stmt

	// While the next tok is not right brace and we are not at the end
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		if p.match(token.END) {
			continue
		}
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

	_, err = p.consume(token.END, "Expect 'end' after block.")
	if err != nil {
		return nil, err
	}

	return stmts, nil
}

func (p *Parser) forStmt() (types.Stmt, error) {
	var err error

	// Dont forget
	// Match advances 'consumes' the next token if matched
	// Check returns wether the next is it or not simply

	var initializer types.Stmt
	if p.match(token.END) { // Just a semicolon, this is directly following the opening (
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

	var condition types.Expr = nil
	if !p.check(token.END) { // See if next token is not semicolon, dont consume it
		if condition, err = p.expression(); err != nil {
			return nil, err // If not, parse expression, not matter what ';' should be at end
		}
	}
	// Consume now
	_, err = p.consume(token.END, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	// Same here but we expect a closing paren instead
	var increment types.Expr = nil
	if !p.check(token.RIGHT_PAREN) {
		if increment, err = p.expression(); err != nil {
			return nil, err
		}
	}

	var body types.Stmt = nil
	if body, err = p.statement(); err != nil {
		return nil, err
	}
	if increment != nil {
		body = types.NewBlock([]types.Stmt{body, types.NewExpression(increment)})
	}
	if condition == nil {
		condition = types.NewLiteralExpr(literal.NewLiteral(true))
	}
	body = types.NewWhile(condition, body)
	if initializer != nil {
		body = types.NewBlock([]types.Stmt{initializer, body})
	}

	return body, nil
}

func (p *Parser) returnStmt() (types.Stmt, error) {
	keyword := p.previous()
	var val types.Expr = nil
	if !p.check(token.END) { // As long as ; is not the next token, cause a ; cant start an expression
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		val = expr
	}
	_, err := p.consume(token.END, "Expect 'end' after return value.")
	if err != nil {
		return nil, err
	}
	return types.NewReturn(keyword, val), nil
}

func (p *Parser) tryStmt() (types.Stmt, error) {
	var woops types.Stmt
	var wert token.Token

	attempt, err := p.statement()
	if err != nil {
		return nil, err
	}

	if p.match(token.WOOPS) {
		wert = p.advance()

		woops, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return types.NewTry(attempt, woops, wert), nil
}

func (p *Parser) whileStmt() (types.Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return types.NewWhile(condition, body), nil
}

func (p *Parser) ifStmt() (types.Stmt, error) {
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var finalBranch types.Stmt
	if p.match(token.ELSE) {
		finalBranch, err = p.statement()
	}

	return types.NewIf(condition, thenBranch, finalBranch), nil
}

func (p *Parser) printStmt() (types.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.END, "Expect 'end' after value.")
	if err != nil {
		return nil, err
	}
	return types.NewPrint(val), nil
}

func (p *Parser) wertStmt() (types.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.END, "Expect 'end' after wert statement.")
	if err != nil {
		return nil, err
	}
	return types.NewWert(val), nil
}

func (p *Parser) exprStmt() (types.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.END, "Expect 'end' after value.")
	if err != nil {
		return nil, err
	}
	return types.NewExpression(val), nil
}

func (p *Parser) expression() (types.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (types.Expr, error) {
	var val types.Expr

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
		if varExpr, ok := expr.(*types.VarExpr); ok {
			name := varExpr.Name
			return types.NewAssignExpr(name, val), nil
		}
	}

	return expr, nil
}

func (p *Parser) or() (types.Expr, error) {
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
		expr = types.NewLogicalExpr(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) and() (types.Expr, error) {
	var err error
	var expr types.Expr

	if expr, err = p.equality(); err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality() // Calling equality then begins seeing the boolean val of the left expression
		if err != nil {
			return nil, err
		}
		expr = types.NewLogicalExpr(expr, operator, right)
	}
	return expr, nil
}

// If we are parsing a == b == c
// We parse a == b, then a == b becomes the left operand of == c, looping
// Returning the entire expression at the end
func (p *Parser) equality() (types.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = types.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (types.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	// While we are currently in a token that is composed of 2 of these
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = types.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (types.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS, token.PLUS_EQUAL, token.MINUS_EQUAL, token.STAR_EQUAL, token.SLASH_EQUAL) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = types.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (types.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = types.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (types.Expr, error) {
	if p.match(token.BANG, token.MINUS) { // If it is ! or -, must be unary
		operator := p.previous()
		right, err := p.unary() // Parse recursively, ie, !!
		if err != nil {
			return nil, err
		}
		return types.NewUnaryExpr(operator, right), nil
	}
	return p.postfix()
	// Must have reached highest level precedence
}

func (p *Parser) postfix() (types.Expr, error) {
	switch p.peekNext().Type {
	case token.PLUS_PLUS, token.MINUS_MINUS:
		left, err := p.call()
		if err != nil {
			return nil, err
		}
		if !p.match(token.PLUS_PLUS, token.MINUS_MINUS) {
			// glorpError.ParserError(p.peek(), )
			return nil, nil
		}
		return types.NewPostfixExpr(left, p.previous()), nil
	}
	return p.call()
}

func (p *Parser) call() (types.Expr, error) {
	var expr types.Expr
	var err error
	if expr, err = p.index(); err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) { // If, after consuming maybe identifier, an opening paren exists
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee types.Expr) (types.Expr, error) {
	var args []types.Expr

	if !p.check(token.RIGHT_PAREN) { // If we dont see right paren as we are walking the args
		for {
			if len(args) >= 255 {
				fmt.Println("Args is over 255!")
				break
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}
	// Finally perform func call
	return types.NewCallExpr(callee, paren, args), nil
}

func (p *Parser) index() (types.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	if expr != nil && p.match(token.LEFT_BRACKET) { // If is non nil expr and leftbracket lies after, could only be index
		idx, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_BRACKET, "Expect ']' to end indexing expression.")
		if err != nil { return nil, err }
		return types.NewIndexExpr(expr, idx), nil
	}
	return expr, nil
}

func (p *Parser) primary() (types.Expr, error) {
	if p.match(token.FALSE) {
		return types.NewLiteralExpr(literal.NewLiteral(false)), nil
	}

	if p.match(token.TRUE) {
		return types.NewLiteralExpr(literal.NewLiteral(true)), nil
	}

	if p.match(token.NIL) {
		return types.NewLiteralExpr(literal.NewLiteral(nil)), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return types.NewLiteralExpr(p.previous().Literal), nil
	}

	// If we see an ident
	if p.match(token.IDENTIFIER) {
		return types.NewVarExpr(p.previous()), nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return types.NewGroupingExpr(expr), nil
	}

	// It has to be in a func that sees if left bracket lies after an expression
	if p.match(token.LEFT_BRACKET) { 
		literalToken := p.previous()
		var data []types.Expr
		for !p.match(token.RIGHT_BRACKET) && !p.isAtEnd() {
			expr, err := p.expression() 
			if err != nil {
				return nil, err
			}
			data = append(data, expr)
			if !p.match(token.COMMA) {
				_, err := p.consume(token.RIGHT_BRACKET, "Expect ']' at the end of glist.")
				if err != nil { return nil, err }
				break
			}
		}
		return types.NewGlistExpr(data, literalToken), nil
	}

	// If we have gotten here, the token given cannot start an expression
	herror.ParserError(p.peek(), "Expect expression.")
	return nil, nil
}

// If passed token is the type of next token, consume it, otherwise error
func (p *Parser) consume(tokType token.TokenType, message string) (token.Token, error) {
	if p.check(tokType) {
		return p.advance(), nil
	} // If next token is passed type, consume it and pass the previous token

	herror.ParserError(p.peek(), message)
	return token.Token{}, fmt.Errorf(message)
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
		case token.CLASS:
		case token.GLUNC:
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

func (p *Parser) matchNext(tokenTypes ...token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	for _, tokType := range tokenTypes {
		if p.checkNext(tokType) {
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

func (p *Parser) checkNext(tokType token.TokenType) bool {
	if p.isAtEnd() { // Because there is no next token
		return false
	}
	return p.peekNext().Type == tokType
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