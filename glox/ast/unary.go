package ast

import "glox/token"

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func NewUnaryExpr(operator token.Token, right Expr) Expr {
	return &UnaryExpr{
		Operator: operator,
		Right:    right,
	}
}

func (u *UnaryExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitUnaryExpr(u)
}
