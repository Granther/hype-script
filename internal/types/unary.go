package types

import (
	"fmt"
	"hype-script/internal/token"
)

type UnaryExpr struct {
	Type     string
	Operator token.Token
	Right    Expr
}

func NewUnaryExpr(operator token.Token, right Expr) Expr {
	return &UnaryExpr{
		Type:     "UnaryExpr",
		Operator: operator,
		Right:    right,
	}
}

func (u *UnaryExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitUnaryExpr(u)
}

func (v *UnaryExpr) GetType() string {
	return v.Type
}

func (v *UnaryExpr) GetVal() string {
	return fmt.Sprintf("%s, %s", v.Operator.String(), v.Right.GetVal())
}
