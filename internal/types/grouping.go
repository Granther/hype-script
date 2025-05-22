package types

import "fmt"

type GroupingExpr struct {
	Type string
	Expr Expr
}

func NewGroupingExpr(expr Expr) Expr {
	return &GroupingExpr{
		Type: "GroupingExpr",
		Expr: expr,
	}
}

func (g *GroupingExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitGroupingExpr(g)
}

func (v *GroupingExpr) GetType() string {
	return v.Type
}

func (v *GroupingExpr) GetVal() string {
	return fmt.Sprintf("%s", v.Expr.GetVal())
}
