package ast

import "glox/token"

type LogicalExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewLogicalExpr(left Expr, operator token.Token, right Expr) Expr {
	return &LogicalExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (v *LogicalExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLogicalExpr(v)
}
