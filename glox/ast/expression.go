package ast

type Visitor interface {
	Print(expr Expr) string
	VisitBinaryExpr(expr *BinaryExpr) (any, error)
	VisitLiteralExpr(expr *LiteralExpr) (any, error)
	VisitUnaryExpr(expr *UnaryExpr) (any, error)
	VisitGroupingExpr(expr *GroupingExpr) (any, error)
	VisitVarExpr(expr *VarExpr) (any, error)
	VisitAssignExpr(expr *AssignExpr) (any, error)
	VisitLogicalExpr(expr *LogicalExpr) (any, error)
	VisitWhileExpr(expr *WhileExpr) (any, error)
	VisitCallExpr(expr *CallExpr) (any, error)
	VisitFunExpr(expr *FunExpr) (any, error)
}

type Expr interface {
	Accept(visitor Visitor) (any, error)
}
