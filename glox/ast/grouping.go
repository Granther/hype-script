package ast

type GroupingExpr struct {
	Expr Expr
}

func NewGroupingExpr(expr Expr) Expr {
	return &GroupingExpr{
		Expr: expr,
	}
}

func (g *GroupingExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitGroupingExpr(g)
}