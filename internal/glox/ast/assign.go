package ast

import (
	"glox/token"
)

type AssignExpr struct {
	Name token.Token
	Val  Expr
}

// Tok for var being assigned to, expr for new val
func NewAssignExpr(name token.Token, val Expr) Expr {
	return &AssignExpr{
		Name: name,
		Val:  val,
	}
}

func (v *AssignExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitAssignExpr(v)
}
