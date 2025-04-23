package types

import (
	"glorp/token"
)

type Expression struct {
	Expr Expr
}

type Print struct {
	Expr Expr
}

type Wert struct {
	Val Expr
}

type Var struct {
	Name        token.Token
	Initializer Expr
}

type Block struct {
	Statements []Stmt
}

type If struct {
	Condition Expr
	Then      Stmt
	Final     Stmt
}

type While struct {
	Condition Expr
	Body      Stmt
}

type Fun struct {
	Params      []token.Token
	Name        token.Token
	Body        []Stmt
	Environment Environment
}

type Return struct {
	Keyword token.Token
	Val     Expr
}

type Try struct {
	Attempt  Stmt
	Woops    Stmt
	WoopsTok token.Token
}

type Class struct {
	Name    token.Token
	Methods []Stmt
}

func NewClass(name token.Token, methods []Stmt) Stmt {
	return &Class{
		Name:    name,
		Methods: methods,
	}
}

func NewTry(attempt Stmt, woops Stmt, woopsTok token.Token) Stmt {
	return &Try{
		Attempt:  attempt,
		Woops:    woops,
		WoopsTok: woopsTok,
	}
}

func NewReturn(keyword token.Token, val Expr) Stmt {
	return &Return{
		Keyword: keyword,
		Val:     val,
	}
}

func NewFun(name token.Token, params []token.Token, body []Stmt, env Environment) Stmt {
	return &Fun{
		Params:      params,
		Name:        name,
		Body:        body,
		Environment: env,
	}
}

func NewWhile(condition Expr, body Stmt) Stmt {
	return &While{
		Condition: condition,
		Body:      body,
	}
}

func NewIf(condition Expr, then Stmt, final Stmt) Stmt {
	return &If{
		Condition: condition,
		Then:      then,
		Final:     final,
	}
}

func NewBlock(statements []Stmt) Stmt {
	return &Block{
		Statements: statements,
	}
}

func NewPrint(expr Expr) Stmt {
	return &Print{
		Expr: expr,
	}
}

func NewWert(expr Expr) Stmt {
	return &Wert{
		Val: expr,
	}
}

func NewExpression(expr Expr) Stmt {
	return &Expression{
		Expr: expr,
	}
}

func NewVar(name token.Token, initializer Expr) Stmt {
	return &Var{
		Name:        name,
		Initializer: initializer,
	}
}

func (e *Print) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrintStmt(e)
}

func (e *Expression) Accept(visitor StmtVisitor) error {
	return visitor.VisitExprStmt(e)
}

func (e *Var) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarStmt(e)
}

func (e *Block) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlockStmt(e)
}

func (e *If) Accept(visitor StmtVisitor) error {
	return visitor.VisitIfStmt(e)
}

func (e *While) Accept(visitor StmtVisitor) error {
	return visitor.VisitWhileStmt(e)
}

func (e *Fun) Accept(visitor StmtVisitor) error {
	return visitor.VisitFunStmt(e)
}

func (e *Return) Accept(visitor StmtVisitor) error {
	return visitor.VisitReturnStmt(e)
}

func (e *Wert) Accept(visitor StmtVisitor) error {
	return visitor.VisitWertStmt(e)
}

func (e *Try) Accept(visitor StmtVisitor) error {
	return visitor.VisitTryStmt(e)
}

func (e *Class) Accept(visitor StmtVisitor) error {
	return visitor.VisitClassStmt(e)
}
