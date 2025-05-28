package hypeinterpreter

// import (
// 	"glorp/environment"
// 	"glorp/literal"
// 	"glorp/token"
// 	"glorp/types"
// 	"testing"
// )

// func TestVisitTryStmt(t *testing.T) {
// 	env := environment.NewEnvironment(nil)
// 	interpreter := NewInterpreter(env)

// 	// Create a mock Try statement
// 	attempt := types.NewExpression(types.NewLiteralExpr(literal.NewLiteral("attempt")))
// 	ohshit := types.NewExpression(types.NewLiteralExpr(literal.NewLiteral("ohshit")))
// 	ohshitTok := token.Token{Type: token.OHSHIT, Lexeme: "ohshit", Line: 1}
// 	tryStmt := types.NewTry(attempt, ohshit, ohshitTok)

// 	// Execute the Try statement
// 	err := interpreter.VisitTryStmt(tryStmt)
// 	if err != nil {
// 		t.Errorf("VisitTryStmt failed: %v", err)
// 	}
// }
