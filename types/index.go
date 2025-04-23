package types

type IndexExpr struct {
	Type  string
	Expr  Expr
	Index Expr
}

func NewIndexExpr(expr Expr, index Expr) Expr {
	return &IndexExpr{
		Type:  "IndexExpr",
		Expr:  expr,
		Index: index,
	}
}

func (v *IndexExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitIndexExpr(v)
}

func (v *IndexExpr) GetType() string {
	return v.Type
}
