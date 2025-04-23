package types

import "glorp/token"

type VarExpr struct {
	Type string
	Name token.Token
}

func NewVarExpr(name token.Token) Expr {
	return &VarExpr{
		Type: "VarExpr",
		Name: name,
	}
}

func (v *VarExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitVarExpr(v)
}

func (v *VarExpr) GetType() string {
	return v.Type
}