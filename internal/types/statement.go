package types

import (
	"fmt"
	"hype-script/internal/token"
)

type Expression struct {
	Expr Expr
}

type Print struct {
	Expr Expr
}

type Var struct {
	Name        token.Token
	Initializer Expr
	Global      bool
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
	Environment EnvironmentHandler
}

type Return struct {
	Keyword token.Token
	Val     Expr
}

type Import struct {
	Lang    token.Token
	Imports []*ImportItem
}

type Access struct {
	Name token.Token
	Expr Expr
}

func NewReturn(keyword token.Token, val Expr) Stmt {
	return &Return{
		Keyword: keyword,
		Val:     val,
	}
}

func NewFun(name token.Token, params []token.Token, body []Stmt, env EnvironmentHandler) Stmt {
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

func NewExpression(expr Expr) Stmt {
	return &Expression{
		Expr: expr,
	}
}

func NewVar(name token.Token, initializer Expr, global bool) Stmt {
	return &Var{
		Name:        name,
		Initializer: initializer,
		Global:      global,
	}
}

func NewImport(lang token.Token, imports []*ImportItem) Stmt {
	return &Import{
		Lang:    lang,
		Imports: imports,
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

func (e *Import) Accept(visitor StmtVisitor) error {
	return visitor.VisitImportStmt(e)
}

// String()
func (e *Print) String() string {
	return fmt.Sprintf("Print ~ Type: %s, Val: %s", e.Expr.GetType(), e.Expr.GetVal())
}

func (e *Expression) String() string {
	return fmt.Sprintf("%s, %s", e.Expr.GetType(), e.Expr.GetVal())
}

func (e *Var) String() string {
	return fmt.Sprintf("%s, %s", e.Initializer.GetType(), e.Initializer.GetVal())
}

func (e *Block) String() string {
	//return fmt.Sprintf("%s, %s", e.Expr.GetType(), e.Expr.GetVal())
	return ""
}

func (e *If) String() string {
	return fmt.Sprintf("%s, %s", e.Condition.GetType(), e.Condition.GetVal())
}

func (e *While) String() string {
	return fmt.Sprintf("%s, %s", e.Condition.GetType(), e.Condition.GetVal())
}

func (e *Fun) String() string {
	return fmt.Sprintf("Name: %s", e.Name.String())
}

func (e *Return) String() string {
	return fmt.Sprintf("%s, %s", e.Val.GetType(), e.Val.GetVal())
}

func (e *Import) String() string {
	return ""
}
