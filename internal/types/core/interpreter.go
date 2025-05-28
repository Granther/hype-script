package core

import "hype-script/internal/types"

type InterpreterHandler interface {
	InterpretStmts(stmts []types.Stmt)
	GetHadRuntimeError() bool 
	ExecuteBlock(stmts []types.Stmt, environment types.Environment) error
	GetGlobals() types.Environment
}
