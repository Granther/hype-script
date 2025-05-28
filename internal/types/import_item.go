package types

import (
	"fmt"
	"hype-script/internal/token"
)

type ImportItem struct {
	Alias token.Token
	Val   token.Token
}

func NewImportItem(alias, val token.Token) *ImportItem {
	return &ImportItem{
		Alias: alias,
		Val:   val,
	}
}

func (i *ImportItem) String() string {
	return fmt.Sprintf("ImportItem -> Alias: %s, Val: %s", &i.Alias.Lexeme, &i.Val.Lexeme)
}
