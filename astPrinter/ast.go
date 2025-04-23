package typesPrinter

import (
	"fmt"
	"glorp/types"
	"glorp/utils"
)

type typesPrinter struct{}

func NewtypesPrinter() types.Visitor {
	return &typesPrinter{}
}

func (a *typesPrinter) VisitBinaryExpr(expr *types.BinaryExpr) (any, error) {
	return utils.Parenthesize(a, expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *typesPrinter) VisitUnaryExpr(expr *types.UnaryExpr) (any, error) {
	return utils.Parenthesize(a, expr.Operator.Lexeme, expr.Right)
}

func (a *typesPrinter) VisitGroupingExpr(expr *types.GroupingExpr) (any, error) {
	return utils.Parenthesize(a, "group", expr.Expr)
}

func (a *typesPrinter) VisitLiteralExpr(expr *types.LiteralExpr) (any, error) {
	if expr.Val.Val == nil {
		return "nil", nil
	}
	return expr.Val.String(), nil
}

func (a *typesPrinter) VisitVarExpr(expr *types.VarExpr) (any, error) {
	return utils.Parenthesize(a, "var")
}

func (a *typesPrinter) VisitAssignExpr(expr *types.AssignExpr) (any, error) {
	return utils.Parenthesize(a, "assign", expr.Val)
}

func (a *typesPrinter) VisitLogicalExpr(expr *types.LogicalExpr) (any, error) {
	return utils.Parenthesize(a, "logical", expr.Left, expr.Right)
}

func (a *typesPrinter) VisitWhileExpr(expr *types.WhileExpr) (any, error) {
	return utils.Parenthesize(a, "while")
}

func (a *typesPrinter) VisitCallExpr(expr *types.CallExpr) (any, error) {
	return utils.Parenthesize(a, "call")
}

func (a *typesPrinter) VisitFunExpr(expr *types.FunExpr) (any, error) {
	return utils.Parenthesize(a, "func")
}

func (a *typesPrinter) VisitReturnExpr(expr *types.ReturnExpr) (any, error) {
	return utils.Parenthesize(a, "return")
}

func (a *typesPrinter) VisitPostfixExpr(expr *types.PostfixExpr) (any, error) {
	return utils.Parenthesize(a, "postfix")
}

func (a *typesPrinter) VisitGlistExpr(expr *types.GlistExpr) (any, error) {
	return utils.Parenthesize(a, "glist")
}

func (a *typesPrinter) Print(expr types.Expr) string {
	val, err := expr.Accept(a)
	if err != nil {
	}
	return fmt.Sprintf("%v", val)
}