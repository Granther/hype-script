package parser

import (
	"errors"
	"fmt"
	herror "hype-script/internal/error"
	"hype-script/internal/literal"
	"hype-script/internal/token"
	"hype-script/internal/types"
)

// Called repeatably to parse a series of statments in a program, perfect place to look for panic
func (p *Parser) declaration() (types.Stmt, error) {
	// If current token is var, we are looking at var decl
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	if p.match(token.FUN) {
		return p.funDeclaration()
	}
	// If not, fallback to standard stmt
	return p.statement()
}

func (p *Parser) funDeclaration() (types.Stmt, error) {
	// Consume name here, match already ate 'fun'
	name, err := p.consume(token.IDENTIFIER, "Expect function name")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after function name")
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

	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before function body.")
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
	var global bool
	if p.match(token.KARAT) {
		global = true
	} else if p.match(token.TILDE) {
		global = false
	}

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

	return types.NewVar(name, initializer, global), nil
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

	if p.match(token.PAR) {
		fmt.Println("Got PAR token")
	}

	if p.match(token.HYP) {
		fmt.Println("Got HYP token")
	}

	if p.match(token.IMPORT) {
		return p.importStmt()
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
		// Always eat newline
		p.match(token.END)

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

func (p *Parser) expression() (types.Expr, error) {
	return p.assignment() // Top level item contained in expression to start recursive descent
}

func (p *Parser) assignment() (types.Expr, error) {
	var val types.Expr

	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
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
	if p.match(token.BANG, token.MINUS, token.KARAT, token.TILDE) { // If it is ! or -, must be unary
		operator := p.previous()
		right, err := p.unary() // Parse recursively, ie, !!
		if err != nil {
			return nil, err
		}
		return types.NewUnaryExpr(operator, right), nil
	}
	return p.access()
	// Must have reached highest level precedence
}

// Mwah (may 30 2025)
func (p *Parser) access() (types.Expr, error) {
	if p.peekNext().Type == token.DOT {
		e, err := p.postfix()
		if err != nil {
			return nil, err
		}
		exprs := []types.Expr{e}
		for p.match(token.DOT) {
			e, err = p.postfix()
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, e)
			if p.match(token.END) {
				return types.NewAccessExpr(exprs), nil
			}
		}
	}
	return p.postfix()
}

func (p *Parser) postfix() (types.Expr, error) {
	switch p.peekNext().Type {
	case token.PLUS_PLUS, token.MINUS_MINUS:
		left, err := p.call()
		if err != nil {
			return nil, err
		}
		if !p.match(token.PLUS_PLUS, token.MINUS_MINUS) {
			msg := "Expected ++ or -- postfix."
			herror.ParserError(p.peek(), msg)
			return nil, errors.New(msg)
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

	// What happens when foo(
	// If right paren is not there we assume params
	//
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
		if err != nil {
			return nil, err
		}
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

	if p.match(token.NEWT) {
		return types.NewLiteralExpr(literal.NewLiteral(nil)), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return types.NewLiteralExpr(p.previous().Literal), nil
	}

	// IDENT DOT IDENT
	// How can this be beautiful?
	// We save Access as an array of exprs

	// Only return access if we see end

	// fmt.Println.Go()

	// How can we do this in an intuitive way?
	// We see expr followed by DOT followed by another expr. We are seeing access

	// If we see an ident
	if p.match(token.IDENTIFIER) {
		return types.NewVarExpr(p.previous()), nil
	}

	// fmt.Println.Thing
	// IDENT DOT IDENT DOT IDENT

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
				if err != nil {
					return nil, err
				}
				break
			}
		}
		return types.NewGlistExpr(data, literalToken), nil
	}

	// If we have gotten here, the token given cannot start an expression
	msg := "Expect expression."
	herror.ParserError(p.peek(), msg)
	return nil, errors.New(msg)
}
