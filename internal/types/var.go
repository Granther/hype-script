package types

import (
	"hype-script/internal/token"
)

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

func (v *VarExpr) GetVal() string {
	return v.Name.String()
}