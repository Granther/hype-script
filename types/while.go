package types

type WhileExpr struct {
	Type      string
	Condition Expr
	Body      Stmt
}

func NewWhileExpr(condition Expr, body Stmt) Expr {
	return &WhileExpr{
		Type:      "WhileExpr",
		Condition: condition,
		Body:      body,
	}
}

func (w *WhileExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitWhileExpr(w)
}

func (v *WhileExpr) GetType() string {
	return v.Type
}