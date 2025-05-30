package core

import "hype-script/internal/types"

type InterpreterHandler interface {
	InterpretStmts(stmts []types.Stmt) error
	GetHadRuntimeError() bool
	ExecuteBlock(stmts []types.Stmt, environment types.EnvironmentHandler) error
	GetGlobals() types.EnvironmentHandler
}
