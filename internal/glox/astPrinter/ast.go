package astPrinter

import (
	"fmt"
	"glox/ast"
	"glox/utils"
)

type AstPrinter struct{}

func NewAstPrinter() ast.Visitor {
	return &AstPrinter{}
}

func (a *AstPrinter) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	return utils.Parenthesize(a, expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	return utils.Parenthesize(a, expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return utils.Parenthesize(a, "group", expr.Expr)
}

func (a *AstPrinter) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	if expr.Val.Val == nil {
		return "nil", nil
	}
	return expr.Val.String(), nil
}

func (a *AstPrinter) VisitVarExpr(expr *ast.VarExpr) (any, error) {
	return utils.Parenthesize(a, "var")
}

func (a *AstPrinter) VisitAssignExpr(expr *ast.AssignExpr) (any, error) {
	return utils.Parenthesize(a, "assign", expr.Val)
}

func (a *AstPrinter) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	return utils.Parenthesize(a, "logical", expr.Left, expr.Right)
}

func (a *AstPrinter) VisitWhileExpr(expr *ast.WhileExpr) (any, error) {
	return utils.Parenthesize(a, "while")
}

func (a *AstPrinter) VisitCallExpr(expr *ast.CallExpr) (any, error) {
	return utils.Parenthesize(a, "call")
}

func (a *AstPrinter) VisitFunExpr(expr *ast.FunExpr) (any, error) {
	return utils.Parenthesize(a, "func")
}

func (a *AstPrinter) Print(expr ast.Expr) string {
	val, err := expr.Accept(a)
	if err != nil {
	}
	return fmt.Sprintf("%v", val)
}
