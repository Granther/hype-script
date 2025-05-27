package types

import "hype-script/internal/token"

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
