package ast

type StmtVisitor interface {
	VisitExprStmt(stmt *Expression) error
	VisitPrintStmt(stmt *Print) error
	VisitVarStmt(stmt *Var) error
	VisitBlockStmt(stmt *Block) error
	VisitIfStmt(stmt *If) error
	VisitWhileStmt(stmt *While) error
	VisitFunStmt(stmt *Fun) error
}