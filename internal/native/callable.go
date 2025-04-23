package native

import (
	"fmt"
	herror "hype-script/internal/error"
	"hype-script/internal/environment"
	"hype-script/internal/types"
)

type Callable interface {
	Arity() int
	Call(interpreter types.Interpreter, args []any) (any, error)
	String() string
}

type GlorpFunction struct {
	Declaration types.Fun
}

func NewGlorpFunction(declaration types.Fun) Callable {
	return &GlorpFunction{
		Declaration: declaration,
	}
}

// Each function gets its own environment to store local vars
// A new environment is necassary when thinking about recursive funs
// They do not share local vars
func (f *GlorpFunction) Call(interpreter types.Interpreter, args []any) (any, error) {
	environment := environment.NewEnvironment(interpreter.GetGlobals())
	for i := 0; i < len(f.Declaration.Params); i++ {
		// Place passed args as accessible in the body locally
		environment.Define(f.Declaration.Params[i].Lexeme, args[i])
	}
	// Call function and discard environ, reverting to prev
	err := interpreter.ExecuteBlock(f.Declaration.Body, environment)
	ret, ok := err.(*herror.ReturnErr)
	if ok {
		return ret.Val, nil
	}

	// _, ok = err.(*glorpError.WertErr) // If it is a wert, allow it up 
	// if ok {
	// 	return nil, err
	// }
	return nil, err
}

func (f *GlorpFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *GlorpFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}
