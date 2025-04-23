package types

import (
	"hype-script/internal/token"
)

type FunExpr struct {
	Type   string
	Params []token.Token
	Name   token.Token
	Body   []Stmt
}

func NewFunExpr(name token.Token, params []token.Token, body []Stmt) Expr {
	return &FunExpr{
		Type:   "FunExpr",
		Params: params,
		Name:   name,
		Body:   body,
	}
}

func (f *FunExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitFunExpr(f)
}

func (v *FunExpr) GetType() string {
	return v.Type
}