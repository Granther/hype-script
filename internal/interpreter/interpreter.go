package interpreter

import (
	"fmt"
	herror "hype-script/internal/error"
	"hype-script/internal/glorpups"
	"hype-script/internal/token"
	"hype-script/internal/types"
	"hype-script/internal/utils"
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

func (i *Interpreter) Interpret(stmts []types.Stmt) {
	for _, stmt := range stmts {
		i.execute(stmt)
		// switch stmt.(type) {
		// case *types.Var:
		// 	i.execute(stmt)
		// case *types.Fun:
		// 	i.execute(stmt)
		// }
	}

	// g, err := i.Environment.Get("mlorp")
	// if err != nil {
	// 	herror.InterpreterRuntimeError(token.Token{}, "mlorp entry glunction not found.")
	// 	return
	// }

	// f, ok := g.(native.Callable)
	// if !ok {
	// 	herror.InterpreterRuntimeError(token.Token{}, "unable to read mlorp entry glunc to callable.")
	// 	return
	// }

	// _, err = f.Call(i, []any{})
	// if err != nil {
	// 	glorpups.InterpreterRuntimeError("uncaught wert arrived in global scope", err)
	// 	return
	// }
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
		herror.InterpreterRuntimeError(operator, "Operand must be number.")
		return fmt.Errorf("unable to convert operand to int")
	}
	return nil
}

func checkNumberOperands(operator token.Token, left any, right any) (float64, float64, error) {
	l, r, ok := utils.ConvFloat(left, right)
	if !ok {
		herror.InterpreterRuntimeError(operator, "Operands must be numbers.")
		return -1, -1, fmt.Errorf("unable to convert operands to int")
	}
	return l, r, nil
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
