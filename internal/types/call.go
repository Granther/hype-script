package types

import (
	"hype-script/internal/token"
)

type CallExpr struct {
	Type   string
	Callee Expr
	Paren  token.Token // Token for closing parens
	Args   []Expr
}

func NewCallExpr(callee Expr, paren token.Token, args []Expr) Expr {
	return &CallExpr{
		Type:   "CallExpr",
		Callee: callee,
		Paren:  paren,
		Args:   args,
	}
}

func (c *CallExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitCallExpr(c)
}

func (v *CallExpr) GetType() string {
	return v.Type
}

func (v *CallExpr) GetVal() string {
	return v.Callee.GetVal()
}
