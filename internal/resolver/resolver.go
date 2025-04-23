package resolver

import (
	"hype-script/internal/interpreter"
	"hype-script/internal/types"
)

type Resolver struct {
	Interprter interpreter.Interpreter
}

func NewResolver(interpreter interpreter.Interpreter) *Resolver {
	return &Resolver{
		Interprter: interpreter,
	}
}

func (r *Resolver) VisitBlockStmt(stmt *types.Block) error {
	// r.beginScope()
	r.resolve(stmt.Statements)
	// r.endScope()
	return nil
}

func (r *Resolver) resolve(stmts []types.Stmt) error {
	for _, stmt := range stmts {
		err := r.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt types.Stmt) error {
	// return stmt.Accept(r)
	return nil
}

func (r *Resolver) resolveExpr(expr types.Expr) error {
	// _, err := expr.Accept(r)
	// if err != nil {
	// 	return err
	// }
	return nil
}