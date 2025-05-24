package parser

import (
	"hype-script/internal/literal"
	"hype-script/internal/token"
	"hype-script/internal/types"
)

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

// func (p *Parser) tryStmt() (types.Stmt, error) {
// 	var woops types.Stmt
// 	var wert token.Token

// 	attempt, err := p.statement()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// if p.match(token.WOOPS) {
// 	// 	wert = p.advance()

// 	// 	woops, err = p.statement()
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// }

// 	return types.NewTry(attempt, woops, wert), nil
// }

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

// func (p *Parser) wertStmt() (types.Stmt, error) {
// 	val, err := p.expression()
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = p.consume(token.END, "Expect 'end' after wert statement.")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return types.NewWert(val), nil
// }

func (p *Parser) exprStmt() (types.Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	// Eat end after expr if exists
	p.match(token.END)
	// _, err = p.consume(token.END, "Expect 'end' after value.")
	// if err != nil {
	// 	return nil, err
	// }
	return types.NewExpression(val), nil
}
