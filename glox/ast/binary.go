package ast

import "glox/token"

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewBinaryExpr(left Expr, operator token.Token, right Expr) Expr {
	return &BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (b *BinaryExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitBinaryExpr(b)
}
