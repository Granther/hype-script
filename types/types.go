package types

import "glorp/token"

type Interpreter interface {
	GetGlobals() Environment
	ExecuteBlock(stmts []Stmt, environment Environment) error
	GetHadRuntimeError() bool
	Interpret(stmts []Stmt)
}

type Parser interface {
	GetHadError() bool
	Parse(tokens []token.Token) []Stmt
}

type Environment interface {
	Get(name string) (any, error)
	Define(name string, val any)
	Assign(name token.Token, val any) error
	String() string
}

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVisitor interface {
	VisitExprStmt(stmt *Expression) error
	VisitPrintStmt(stmt *Print) error
	VisitVarStmt(stmt *Var) error
	VisitBlockStmt(stmt *Block) error
	VisitIfStmt(stmt *If) error
	VisitWhileStmt(stmt *While) error
	VisitFunStmt(stmt *Fun) error
	VisitReturnStmt(stmt *Return) error
	VisitWertStmt(stmt *Wert) error
	VisitTryStmt(stmt *Try) error
	VisitClassStmt(stmt *Class) error
}

type Visitor interface {
	Print(expr Expr) (string, error)
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
	VisitReturnExpr(expr *ReturnExpr) (any, error)
	VisitPostfixExpr(expr *PostfixExpr) (any, error)
	VisitGlistExpr(expr *GlistExpr) (any, error)
	VisitIndexExpr(expr *IndexExpr) (any, error)
}

type Expr interface {
	Accept(visitor Visitor) (any, error)
	GetType() string
}
