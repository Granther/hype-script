package types

import (
	"fmt"
	"hype-script/internal/token"
)

type AssignExpr struct {
	Type string
	Name token.Token
	Val  Expr
}

// Tok for var being assigned to, expr for new val
func NewAssignExpr(name token.Token, val Expr) Expr {
	return &AssignExpr{
		Type: "AssignExpr",
		Name: name,
		Val:  val,
	}
}

func (v *AssignExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitAssignExpr(v)
}

func (v *AssignExpr) GetType() string {
	return v.Type
}

func (v *AssignExpr) GetVal() string {
	return fmt.Sprintf("%s, %s", v.Name.String(), v.Val.GetVal())
}
