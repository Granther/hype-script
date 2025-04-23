package ast

import "glox/token"

type VarExpr struct {
	Name token.Token
}

func NewVarExpr(name token.Token) Expr {
	return &VarExpr{
		Name: name,
	}
}

func (v *VarExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitVarExpr(v)
}
