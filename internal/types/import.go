package types

import (
	"fmt"
	"hype-script/internal/token"
)

type ImportExpr struct {
	Type    string
	Keyword token.Token
	Val     Expr
}

func NewImportExpr(keyword token.Token, val Expr) Expr {
	return &ImportExpr{
		Type:    "ImportExpr",
		Keyword: keyword,
		Val:     val,
	}
}

func (r *ImportExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitImportExpr(r)
}

func (v *ImportExpr) GetType() string {
	return v.Type
}

func (v *ImportExpr) GetVal() string {
	return fmt.Sprintf("%s, %s", v.Keyword.String(), v.Val.GetVal())
}
