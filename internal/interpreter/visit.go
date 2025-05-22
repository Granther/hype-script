package interpreter

import (
	"fmt"
	"hype-script/internal/environment"
	herror "hype-script/internal/error"
	"hype-script/internal/native"
	"hype-script/internal/token"
	"hype-script/internal/types"
	"hype-script/internal/utils"
	"reflect"
)

func (i *Interpreter) VisitBinaryExpr(expr *types.BinaryExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.MINUS:
		l, r, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l - r, nil
	case token.SLASH:
		l, r, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l / r, nil
	case token.STAR:
		l, r, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l * r, nil
	case token.GREATER:
		l, r, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l > r, nil
	case token.GREATER_EQUAL:
		l, r, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l >= r, nil
	case token.LESS:
		l, r, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l < r, nil
	case token.LESS_EQUAL:
		l, r, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l <= r, nil
	case token.BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case token.PLUS_EQUAL, token.MINUS_EQUAL, token.STAR_EQUAL, token.SLASH_EQUAL:
		l, r, ok := utils.ConvFloat(left, right) // See if it is int
		if !ok {
			break
		}
		var val any
		switch expr.Operator.Type {
		case token.PLUS_EQUAL:
			val = l + r
		case token.MINUS_EQUAL:
			val = l - r
		case token.STAR_EQUAL:
			val = l * r
		case token.SLASH_EQUAL:
			val = l / r
		}
		if err := i.postfixAssign(expr.Left, val); err != nil { // Attempt to assign to var if one exists
			return nil, err
		}
		return val, nil
	case token.PLUS:
		l, r, ok := utils.ConvFloat(left, right) // See if it is int
		if ok {
			return l + r, nil
		} else if reflect.TypeOf(left).Kind().String() == "string" && reflect.TypeOf(right).Kind().String() == "string" {
			return fmt.Sprintf("%v", left) + fmt.Sprintf("%v", right), nil // If they are both strings then concat
		}
	}

	return utils.Parenthesize(i, expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (i *Interpreter) postfixAssign(expr types.Expr, val any) error {
	variable, ok := expr.(*types.VarExpr)
	if ok {
		return i.Environment.Assign(variable.Name, val)
	}
	return nil
}

func (i *Interpreter) VisitUnaryExpr(expr *types.UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.BANG:
		return !i.isTruthy(right), nil // If the expression on the right is
	case token.MINUS:
		val, ok := right.(float64)
		if ok {
			return -val, nil
		}
	}

	return utils.Parenthesize(i, expr.Operator.Lexeme, expr.Right)
}

func (i *Interpreter) VisitPostfixExpr(expr *types.PostfixExpr) (any, error) {
	var val float64
	var ok bool

	// Find actual value of value of expr to perform oper on (i in i++)
	left, err := i.evaluate(expr.Val)
	if err != nil {
		return nil, err
	}

	// Do operation
	switch expr.Operator.Type {
	case token.PLUS_PLUS:
		val, ok = left.(float64)
		if ok {
			val = val + 1
		}
	case token.MINUS_MINUS:
		val, ok = left.(float64)
		if ok {
			val = val - 1
		}
	}

	// If expr is a variable, reassign the variable
	variable, ok := expr.Val.(*types.VarExpr)
	if ok {
		if err = i.Environment.Assign(variable.Name, val); err != nil {
			return nil, err
		}
	}

	return utils.Parenthesize(i, expr.Operator.Lexeme, expr.Val)
}

// Recursively looks through layered parens
func (i *Interpreter) VisitGroupingExpr(expr *types.GroupingExpr) (any, error) {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitLiteralExpr(expr *types.LiteralExpr) (any, error) {
	return expr.Val.Val, nil
}

func (i *Interpreter) VisitWhileExpr(expr *types.WhileExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitReturnExpr(expr *types.ReturnExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitCallExpr(expr *types.CallExpr) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	var args []any
	for _, arg := range expr.Args {
		val, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}

	fun, ok := callee.(native.Callable)
	if !ok {
		herror.InterpreterRuntimeError(expr.Paren, fmt.Sprintf("Expected identifier, got type %T", fun))
	}

	// Check that the function has the right amount of args passed, args same len as params
	if len(args) != fun.Arity() {
		herror.InterpreterRuntimeError(expr.Paren, fmt.Sprintf("Expected %d args but got %d.", fun.Arity(), len(args)))
	}

	x, err := fun.Call(i, args)
	if err != nil {
		switch err.(type) {
		case *herror.WertErr:
		}
	}

	return x, err
}

func (i *Interpreter) VisitExprStmt(stmt *types.Expression) error {
	_, err := i.evaluate(stmt.Expr)
	return err
}

func (i *Interpreter) VisitReturnStmt(stmt *types.Return) error {
	if stmt.Val != nil {
		v, err := i.evaluate(stmt.Val)
		if err != nil {
			return err
		}
		return herror.NewReturnErr(v)
	}
	return nil // No return val
}

func (i *Interpreter) VisitPrintStmt(stmt *types.Print) error {
	val, err := i.evaluate(stmt.Expr)
	if err != nil {
		return err
	}
	fmt.Println(utils.Stringify(val))
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *types.Var) error {
	var val any
	var err error
	if stmt.Initializer != nil { // If variable has initializer " = 10"
		val, err = i.evaluate(stmt.Initializer) // Evaluate it
		if err != nil {
			return err
		}
	}
	// A variable without an initializer is declared but NOT assigned, error to access before assignment
	// But here, we auto assign the var to nil
	// Give variable value
	i.Environment.Define(stmt.Name.Lexeme, val)
	return nil
}

func (i *Interpreter) VisitFunStmt(stmt *types.Fun) error {
	// Take fun syntax node
	function := native.NewGlorpFunction(*stmt)
	i.Environment.Define(stmt.Name.Lexeme, function) // Add function to global environment by name, can be used anywhere now
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *types.Block) error {
	return i.ExecuteBlock(stmt.Statements, environment.NewEnvironment(i.Environment))
}

func (i *Interpreter) VisitWhileStmt(stmt *types.While) error {
	for {
		val, err := i.evaluate(stmt.Condition)
		if err != nil {
			return err
		}

		if !i.isTruthy(val) {
			break
		}

		err = i.execute(stmt.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *types.If) error {
	val, err := i.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	// Then is like the first if, if true run that statement
	if i.isTruthy(val) {
		err = i.execute(stmt.Then)
		if err != nil {
			return err
		}
	} else if stmt.Final != nil { // Final (else keyword is taken in Go)
		i.execute(stmt.Final)
	}
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *types.AssignExpr) (any, error) {
	val, err := i.evaluate(expr.Val)
	if err != nil {
		return nil, err
	}
	err = i.Environment.Assign(expr.Name, val)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (i *Interpreter) VisitVarExpr(expr *types.VarExpr) (any, error) {
	return i.Environment.Get(expr.Name.Lexeme)
}

func (i *Interpreter) VisitGlistExpr(expr *types.GlistExpr) (any, error) {
	return expr.Data, nil
}

func (i *Interpreter) VisitIndexExpr(expr *types.IndexExpr) (any, error) {
	indexVal, err := i.evaluateIndex(expr.Index)
	if err != nil {
		return nil, err
	}

	switch expr.Expr.(type) {
	case *types.LiteralExpr:
		return i.indexLiteral(expr.Expr, indexVal)
	case *types.GlistExpr:
		return i.indexGlist(expr.Expr, indexVal)
	case *types.VarExpr:
		return i.indexVar(expr.Expr, indexVal)
	}
	return nil, nil
}

func (i *Interpreter) VisitFunExpr(expr *types.FunExpr) (any, error) {
	return i.Environment.Get(expr.Name.Lexeme)
}

func (i *Interpreter) VisitLogicalExpr(expr *types.LogicalExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == token.OR {
		// Short circut, allowing us to exit expression as true/false before evaling the right operand
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		// If the operator is AND, see if it is not truthy, another early escape
		if !i.isTruthy(left) {
			return left, nil
		}
	}
	// The right operator will decide the expressed value of the expression
	// At this point
	// If or: left is false
	// If and: left is true
	return i.evaluate(expr.Right)
}
