package types

import "glorp/token"

type ReturnExpr struct {
	Type    string
	Keyword token.Token
	Val     Expr
}

func NewReturnExpr(keyword token.Token, val Expr) Expr {
	return &ReturnExpr{
		Type:    "ReturnExpr",
		Keyword: keyword,
		Val:     val,
	}
}

func (r *ReturnExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitReturnExpr(r)
}

func (v *ReturnExpr) GetType() string {
	return v.Type
}