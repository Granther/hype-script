package core

import "hype-script/internal/token"

type ScannerHandler interface {
	ScanTokens(source string) ([]token.Token, error)
}
