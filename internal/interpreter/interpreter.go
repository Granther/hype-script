package interpreter

import (
	"fmt"
	"hype-script/internal/environment"
	glorpError "hype-script/internal/error"
	"hype-script/internal/glorpups"
	"hype-script/internal/native"
	"hype-script/internal/token"
	"hype-script/internal/types"
	"hype-script/internal/utils"
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
	Globals         types.Environment
	Environment     types.Environment
}

func NewInterpreter(env types.Environment) types.Interpreter {
	// globals := environment.NewEnvironment(nil)
	// globals.Define("clock", native.NewClockCallable())
	return &Interpreter{
		// Pass nil because we want this to point to the global scope
		Globals:         env,
		Environment:     env,
		HadRuntimeError: false,
	}
}

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
		glorpError.InterpreterRuntimeError(expr.Paren, fmt.Sprintf("Expected identifier, got type %T", fun))
	}

	// Check that the function has the right amount of args passed, args same len as params
	if len(args) != fun.Arity() {
		glorpError.InterpreterRuntimeError(expr.Paren, fmt.Sprintf("Expected %d args but got %d.", fun.Arity(), len(args)))
	}

	x, err := fun.Call(i, args)
	if err != nil {
		switch err.(type) {
		case *glorpError.WertErr:
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
		return glorpError.NewReturnErr(v)
	}
	return nil // No return val
}

// func (i *Interpreter) VisitWertStmt(stmt *types.Wert) error {
// 	if stmt.Val != nil {
// 		v, err := i.evaluate(stmt.Val)
// 		if err != nil {
// 			return err
// 		}
// 		return glorpError.NewWertErr(v)
// 	}
// 	return nil
// }

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

// func (i *Interpreter) VisitClassStmt(stmt *types.Class) error {
// 	return nil
// }

// func (i *Interpreter) VisitTryStmt(stmt *types.Try) error {
// 	err := i.execute(stmt.Attempt)
// 	switch err.(type) {
// 	case *glorpError.ReturnErr:
// 		return err
// 	case glorpups.Glorpup:
// 		wertVal, _ := err.(*glorpError.WertErr)
// 		newVar := types.NewVar(stmt.WoopsTok, types.NewLiteralExpr(literal.NewLiteral(wertVal)))
// 		block := types.NewBlock([]types.Stmt{newVar, stmt.Woops})
// 		err = i.execute(block)
// 		// case *glorpError.WertErr:
// 		// 	wertVal, _ := err.(*glorpError.WertErr)
// 		// 	newVar := types.NewVar(stmt.WoopsTok, types.NewLiteralExpr(literal.NewLiteral(wertVal)))
// 		// 	block := types.NewBlock([]types.Stmt{newVar, stmt.Woops})
// 		// 	err = i.execute(block)
// 	}
// 	return err
// }

func (i *Interpreter) ExecuteBlock(stmts []types.Stmt, environment types.Environment) error {
	prev := i.Environment // Save old, for setting back later

	// Change to new block and execute from that env
	i.Environment = environment
	for _, stmt := range stmts {
		// If a stmt is wert,
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

func (i *Interpreter) evaluateIndex(index types.Expr) (any, error) {
	switch index.(type) {
	case *types.VarExpr:
		variable, ok := index.(*types.VarExpr)
		if !ok {
			fmt.Println("not good got var in indexexpr")
			return nil, nil
		}
		val, err := i.Environment.Get(variable.Name.Lexeme)
		if err != nil {
			return nil, err
		}
		return val, nil
	case *types.LiteralExpr:
		idxLit, ok := index.(*types.LiteralExpr)
		if ok {
			return idxLit.GetRawVal(), nil
		}
	}
	return nil, nil
}

func (i *Interpreter) indexGlist(expr types.Expr, index any) (any, error) {
	if float, ok := index.(float64); ok {
		idx := int(float)
		glist := expr.(*types.GlistExpr)
		if idx > len(glist.Data)-1 {
			return nil, glorpups.NewIndexBoundsGlorpup(glist.GetToken(), "Index out of bounds", nil)
		}
		return i.evaluate(glist.Data[idx])
	}
	return nil, glorpups.NewIndexBoundsGlorpup(token.Token{}, "Incorrect type for indexing Glist.", nil)
}

func (i *Interpreter) indexVar(expr types.Expr, index any) (any, error) {
	var ok bool

	variable, ok := expr.(*types.VarExpr) // Collapse to variable expr
	if ok {
		varVal, err := i.Environment.Get(variable.Name.Lexeme) // Get var val from env
		if err != nil {
			return nil, err
		}

		exprList, ok := varVal.([]types.Expr) // Turn into slice of exprs
		if ok {
			indexedVal := exprList[int(index.(float64))]
			v, ok := indexedVal.(*types.LiteralExpr)
			if ok {
				return v.Val.Val, nil
			}
		}
	}
	return nil, fmt.Errorf("unable to index variable expression")
}

func (i *Interpreter) indexLiteral(expr types.Expr, index any) (any, error) {
	var str string
	var ok bool

	val := expr.(*types.LiteralExpr)
	str, ok = val.GetRawVal().(string)
	if !ok {
		return nil, fmt.Errorf("literal is not iterable")
	}

	switch index.(type) {
	case int:
		return str[index.(int)], nil
	case string:
		findChar := index.(string)
		for i, r := range str {
			if string(r) == findChar {
				return i, nil
			}
		}
		return nil, fmt.Errorf("char not in string")
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

func (i *Interpreter) Print(expr types.Expr) (string, error) {
	val, err := expr.Accept(i)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", val), nil
}

// Calls the visit method for whatever dtype it is
func (i *Interpreter) evaluate(expr types.Expr) (any, error) {
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
		glorpError.InterpreterRuntimeError(operator, "Operand must be number.")
		return fmt.Errorf("unable to convert operand to int")
	}
	return nil
}

func checkNumberOperands(operator token.Token, left any, right any) (float64, float64, error) {
	l, r, ok := utils.ConvFloat(left, right)
	if !ok {
		glorpError.InterpreterRuntimeError(operator, "Operands must be numbers.")
		return -1, -1, fmt.Errorf("unable to convert operands to int")
	}
	return l, r, nil
}

func (i *Interpreter) Interpret(stmts []types.Stmt) {
	for _, stmt := range stmts {
		switch stmt.(type) {
		case *types.Var:
			i.execute(stmt)
		case *types.Fun:
			i.execute(stmt)
		}
	}

	g, err := i.Environment.Get("mlorp")
	if err != nil {
		glorpError.InterpreterRuntimeError(token.Token{}, "mlorp entry glunction not found.")
		return
	}

	f, ok := g.(native.Callable)
	if !ok {
		glorpError.InterpreterRuntimeError(token.Token{}, "unable to read mlorp entry glunc to callable.")
		return
	}

	_, err = f.Call(i, []any{})
	if err != nil {
		glorpups.InterpreterRuntimeError("uncaught wert arrived in global scope", err)
		return
	}
}

func (i *Interpreter) execute(stmt types.Stmt) error {
	return stmt.Accept(i)
}

func (i *Interpreter) GetGlobals() types.Environment {
	return i.Environment
}

func (i *Interpreter) GetHadRuntimeError() bool {
	return i.HadRuntimeError
}

// Variable declarations are statements, because we are doing something
// Now we add a declaration grammar rule to out syntax
// Allow for declaring var, funcs, classes
// Can fall through to a statement

// Scoping
// Create an entirely new environment inside each scope block
// We can discard this entire env and not be afraid of deleting global vars of the same name
// Shadowing
// When 2 vars (maybe global and local) have the same name. The local var ctypess a
// 'shadow' over the global one, hiding it
// Environment chaining
// Each environment has link to the env above, all ending in the global scope
// We walk up this chain when looks for vars

// Blocks
// A possibly empty series of statements for decls in curly braces
// A block is a statement, can appear anywhere a statement is allowed

// If we find a return statement, go up to main

// OLD

// x := types.NewCallExpr(f, token.RIGHT_PAREN, []types.Expr{})
// x.Accept(i)
// function := native.NewGlorpFunction(x)
// function.Call(i, []any{})

// i.execute(g.(types.Stmt))
// glorpFunc, ok := g.(native.Callable)
// if !ok {
// 	glorpError.InterpreterRuntimeError(token.Token{}, "unable to convert mlorp to statement.")
// 	return
// }

// for _, stmt := range stmts {
// 	if i.execute(stmt) != nil {
// 		fmt.Println("Error in interpret")
// 		i.HadRuntimeError = true
// 		return
// 	}
// }
// if i.execute(glorpFunc) != nil {
// 	fmt.Println("Error in interpret")
// 	i.HadRuntimeError = true
// 	return
// }
// glorpFunc.Call(i, []any{})

// for _, stmt := range stmts {
// 	if i.execute(stmt) != nil {
// 		fmt.Println("Error in interpret")
// 		i.HadRuntimeError = true
// 		return
// 	}
// }

// Second pass: execute normal statements
// for _, stmt := range stmts {
//     if _, ok := stmt.(*types.Var); !ok {
//         i.execute(stmt)
//     }
// }
