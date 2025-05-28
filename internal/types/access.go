package types

import (
	"fmt"
	"hype-script/internal/token"
)

type AccessExpr struct {
	Type string
	Name token.Token
	Expr Expr
}

func NewAccessExpr(name token.Token, expr Expr) Expr {
	return &AccessExpr{
		Type: "AccessExpr",
		Name: name,
		Expr: expr,
	}
}

func (b *AccessExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitAccessExpr(b)
}

func (v *AccessExpr) GetType() string {
	return v.Type
}

func (v *AccessExpr) GetVal() string {
	return fmt.Sprintf("%s, %s", v.Name.String(), v.Expr.GetVal())
}
