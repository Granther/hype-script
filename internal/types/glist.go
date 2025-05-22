package types

import (
	"fmt"
	"hype-script/internal/token"
)

type GlistExpr struct {
	Type  string
	Token token.Token
	Data  []Expr
}

func NewGlistExpr(data []Expr, token token.Token) Expr {
	return &GlistExpr{
		Type:  "GlistExpr",
		Token: token,
		Data:  data,
	}
}

func (v *GlistExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitGlistExpr(v)
}

func (v *GlistExpr) GetType() string {
	return v.Type
}

func (v *GlistExpr) GetToken() token.Token {
	return v.Token
}

func (v *GlistExpr) GetVal() string {
	return fmt.Sprintf("%s", v.Token.String()) // Wasnt going to deal with parsing list
}
