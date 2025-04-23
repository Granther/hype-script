package main

import (
	"fmt"
	"glox/ast"
	"glox/environment"
	"glox/token"
	"glox/utils"
	"reflect"
)

// State and statements
// Since a statement does not evaluate to a value it must something
// This is called a side effect
// Modify some internal state, ie, save, free
// Print something for the user

// A program is a arbitrarily long list of statements preceding an EOF token
// A Statement can either be an expression or print
// An expression statement ends in a ;
// A print statement is the same but begins with "print"

// Operands of + are always an expression, the HAVE to have a value
// The body of a while loop is a statement, but, that statement can BE an expression

type Interpreter struct {
	HadRuntimeError bool
	Globals         *environment.Environment
	Environment     *environment.Environment
}

func NewInterpreter() *Interpreter {
	globals := environment.NewEnvironment(nil)
	globals.Define("clock", NewClockCallable())
	return &Interpreter{
		// Pass nil because we want this to point to the global scope
		Globals:         globals,
		Environment:     globals,
		HadRuntimeError: false,
	}
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
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
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	case token.PLUS:
		l, r, ok := utils.ConvFloat(left, right) // See if it is int
		if ok {
			return l + r, nil
		} else if reflect.TypeOf(l).Kind().String() == "string" && reflect.TypeOf(r).Kind().String() == "string" {
			return l + r, nil // If they are both strings then concat
		}
	}

	return utils.Parenthesize(i, expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.BANG:
		return !i.isTruthy(right), nil // If the expression on the right is
	case token.MINUS:
		val, ok := right.(int)
		if ok {
			return -val, nil
		}
	}

	return utils.Parenthesize(i, expr.Operator.Lexeme, expr.Right)
}

// Recursively looks through layered parens
func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	return expr.Val.Val, nil
}

func (i *Interpreter) VisitWhileExpr(expr *ast.WhileExpr) (any, error) {
	return nil, nil
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) (any, error) {
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

	fun, ok := callee.(Callable)
	if !ok {
		InterpreterRuntimeError(expr.Paren, fmt.Sprintf("Expected identifier, got type %T", fun))
	}

	// Check that the function has the right amount of args passed, args same len as params
	if len(args) != fun.Arity() {
		InterpreterRuntimeError(expr.Paren, fmt.Sprintf("Expected %d args but got %d.", fun.Arity(), len(args)))
	}

	return fun.Call(i, args)
}

func (i *Interpreter) VisitExprStmt(stmt *ast.Expression) error {
	i.evaluate(stmt.Expr)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.Print) error {
	val, _ := i.evaluate(stmt.Expr)
	fmt.Println(utils.Stringify(val))
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *ast.Var) error {
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

func (i *Interpreter) VisitFunStmt(stmt *ast.Fun) error {
	// Take fun syntax node
	function := NewGloxFunction(*stmt)
	i.Environment.Define(stmt.Name.Lexeme, function) // Add function to global environment by name, can be used anywhere now
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.Block) error {
	return i.executeBlock(stmt.Statements, environment.NewEnvironment(i.Environment))
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.While) error {
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

func (i *Interpreter) VisitIfStmt(stmt *ast.If) error {
	val, err := i.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	// Then is like the first if, if true run that statement
	if i.isTruthy(val) {
		i.execute(stmt.Then)
	} else if stmt.Final != nil { // Final (else keyword is taken in Go)
		i.execute(stmt.Final)
	}
	return nil
}

func (i *Interpreter) executeBlock(stmts []ast.Stmt, environment *environment.Environment) error {
	prev := i.Environment // Save old, for setting back later

	// Change to new block and execute from that env
	i.Environment = environment
	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}

	// Always change back to original env
	end := func() {
		i.Environment = prev
	}
	defer end()

	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) (any, error) {
	val, err := i.evaluate(expr.Val)
	if err != nil {
		return nil, err
	}
	i.Environment.Assign(expr.Name, val)
	return val, nil
}

func (i *Interpreter) VisitVarExpr(expr *ast.VarExpr) (any, error) {
	return i.Environment.Get(expr.Name)
}

func (i *Interpreter) VisitFunExpr(expr *ast.FunExpr) (any, error) {
	return i.Environment.Get(expr.Name)
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
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

func (i *Interpreter) Print(expr ast.Expr) string {
	val, err := expr.Accept(i)
	if err != nil {
	}
	return fmt.Sprintf("%v", val)
}

// Calls the visit method for whatever dtype it is
func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

// False if something is falsey, bool val if val passed is bool, like nil, true for everyhting else
func (i *Interpreter) isTruthy(val any) bool {
	if val == nil { // If nil, obvo falsey
		return false
	}
	b, ok := val.(bool) // If type is bool, return the bool
	if ok {
		return b
	}
	return true // For everything except nil and false
}

func (i *Interpreter) isEqual(a, b any) bool {
	return a == b
}

func checkNumberOperand(operator token.Token, operand any) error {
	_, ok := utils.IsFloat(operand)

	if !ok {
		InterpreterRuntimeError(operator, "Operand must be number.")
		return fmt.Errorf("unable to convert operand to int")
	}
	return nil
}

func checkNumberOperands(operator token.Token, left any, right any) (float64, float64, error) {
	l, r, ok := utils.ConvFloat(left, right)
	if !ok {
		InterpreterRuntimeError(operator, "Operands must be numbers.")
		return -1, -1, fmt.Errorf("unable to convert operands to int")
	}
	return l, r, nil
}

func (i *Interpreter) interpret(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		if i.execute(stmt) != nil {
			fmt.Println("Error in interpret")
			i.HadRuntimeError = true
			return
		}
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) error {
	return stmt.Accept(i)
}

// Variable declarations are statements, because we are doing something
// Now we add a declaration grammar rule to out syntax
// Allow for declaring var, funcs, classes
// Can fall through to a statement

// Scoping
// Create an entirely new environment inside each scope block
// We can discard this entire env and not be afraid of deleting global vars of the same name
// Shadowing
// When 2 vars (maybe global and local) have the same name. The local var casts a
// 'shadow' over the global one, hiding it
// Environment chaining
// Each environment has link to the env above, all ending in the global scope
// We walk up this chain when looks for vars

// Blocks
// A possibly empty series of statements for decls in curly braces
// A block is a statement, can appear anywhere a statement is allowed
