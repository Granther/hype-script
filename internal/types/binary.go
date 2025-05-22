package types

import (
	"fmt"
	"hype-script/internal/token"
)

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

func (v *BinaryExpr) GetVal() string {
	return fmt.Sprintf("%s, %s, %s", v.Left.GetVal(), v.Operator.String(), v.Right.GetVal())
}
