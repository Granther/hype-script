package core

import (
	"hype-script/internal/token"
	"hype-script/internal/types"
)

type ParserHandler interface {
	ParseTokens(tokens []token.Token) ([]types.Stmt, error)
	GetHadError() bool 
}
