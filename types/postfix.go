package types

import (
	"glorp/token"
)

type PostfixExpr struct {
	Type     string
	Val      Expr
	Operator token.Token
}

// Tok for var being assigned to, expr for new val
func NewPostfixExpr(val Expr, operator token.Token) Expr {
	return &PostfixExpr{
		Type:     "PostfixExpr",
		Val:      val,
		Operator: operator,
	}
}

func (v *PostfixExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitPostfixExpr(v)
}

func (v *PostfixExpr) GetType() string {
	return v.Type
}