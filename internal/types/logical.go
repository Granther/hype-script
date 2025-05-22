package types

import (
	"fmt"
	"hype-script/internal/token"
)

type LogicalExpr struct {
	Type     string
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewLogicalExpr(left Expr, operator token.Token, right Expr) Expr {
	return &LogicalExpr{
		Type:     "LogicalExpr",
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (v *LogicalExpr) Accept(visitor Visitor) (any, error) {
	return visitor.VisitLogicalExpr(v)
}

func (v *LogicalExpr) GetType() string {
	return v.Type
}

func (v *LogicalExpr) GetVal() string {
	return fmt.Sprintf("%s, %s, %s", v.Left.GetVal(), v.Operator.String(), v.Right.GetVal())
}
