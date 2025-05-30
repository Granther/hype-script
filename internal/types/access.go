package types

type AccessExpr struct {
	Type  string
	Exprs []Expr
}

func NewAccessExpr(exprs []Expr) Expr {
	return &AccessExpr{
		Type:  "AccessExpr",
		Exprs: exprs,
	}
}

func (b *AccessExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitAccessExpr(b)
}

func (v *AccessExpr) GetType() string {
	return v.Type
}

func (v *AccessExpr) GetVal() string {
	return "Not impl"
}
