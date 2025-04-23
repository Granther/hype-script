package ast

import "glox/token"

type CallExpr struct {
	Callee Expr
	Paren  token.Token // Token for closing parens
	Args   []Expr
}

func NewCallExpr(callee Expr, paren token.Token, args []Expr) Expr {
	return &CallExpr{
		Callee: callee,
		Paren:  paren,
		Args:   args,
	}
}

func (c *CallExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitCallExpr(c)
}
