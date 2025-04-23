package ast

import "glox/token"

type FunExpr struct {
	Params []token.Token
	Name   token.Token
	Body   []Stmt
}

func NewFunExpr(name token.Token, params []token.Token, body []Stmt) Expr {
	return &FunExpr{
		Params: params,
		Name: name,
		Body: body,
	}
}

func (f *FunExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitFunExpr(f)
}
