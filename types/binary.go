package types

import "glorp/token"

type BinaryExpr struct {
	Type     string
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewBinaryExpr(left Expr, operator token.Token, right Expr) Expr {
	return &BinaryExpr{
		Type:     "BinaryExpr",
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (b *BinaryExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitBinaryExpr(b)
}

func (v *BinaryExpr) GetType() string {
	return v.Type
}