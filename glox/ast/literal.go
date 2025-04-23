package ast

import "glox/literal"

type LiteralExpr struct {
	Val *literal.Literal
}

func NewLiteralExpr(val *literal.Literal) Expr {
	return &LiteralExpr{
		Val: val,
	}
}

func (l *LiteralExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLiteralExpr(l)
}
