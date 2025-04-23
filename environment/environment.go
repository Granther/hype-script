package environment

import (
	"fmt"
	"glorp/token"
	"glorp/types"
)

type Environment struct {
	Enlcosing types.Environment
	Values    map[string]any
}

func NewEnvironment(enclosing types.Environment) *Environment {
	return &Environment{
		Enlcosing: enclosing,
		Values:    make(map[string]any),
	}
}

// We want undefined vars to be bugs
func (e *Environment) Get(name string) (any, error) {
	val, ok := e.Values[name]
	if ok {
		return val, nil
	}

	// Start a recursive chain to call up to higher envs, looking for the var
	if e.Enlcosing != nil {
		return e.Enlcosing.Get(name)
	}

	return nil, fmt.Errorf("undefined variable %s", name)
}

func (e *Environment) Define(name string, val any) {
	e.Values[name] = val
}

// Cannot create new var when assigning, thus runtime error
// Has to be runtime because what if a condition must be met to create a var/access it
func (e *Environment) Assign(name token.Token, val any) error {
	_, ok := e.Values[name.Lexeme]
	if ok {
		e.Values[name.Lexeme] = val
		return nil
	}

	// If the name in not in the local scope, check the one above and so on
	if e.Enlcosing != nil {
		e.Enlcosing.Assign(name, val)
		return nil
	}

	return fmt.Errorf("undefined variable %s", name.Lexeme)
}

func (e *Environment) String() string {
	return fmt.Sprintf("%v", e.Values)
}
