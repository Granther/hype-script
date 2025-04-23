package ast

type WhileExpr struct {
	Condition Expr
	Body      Stmt
}

func NewWhileExpr(condition Expr, body Stmt) Expr {
	return &WhileExpr{
		Condition: condition,
		Body:      body,
	}
}

func (w *WhileExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitWhileExpr(w)
}
