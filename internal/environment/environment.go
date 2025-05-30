package environment

import (
	"fmt"
	"hype-script/internal/types"
)

type Environment struct {
	Enlcosing types.EnvironmentHandler
	Values    map[string]any
}

func NewEnvironment(enclosing types.EnvironmentHandler) *Environment {
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
func (e *Environment) Assign(name string, val any) error {
	_, ok := e.Values[name]
	if ok {
		e.Values[name] = val
		return nil
	}

	// If the name in not in the local scope, check the one above and so on
	if e.Enlcosing != nil {
		e.Enlcosing.Assign(name, val)
		return nil
	}

	return fmt.Errorf("undefined variable %s", name)
}

func (e *Environment) String() string {
	return fmt.Sprintf("%v", e.Values)
}

func (e *Environment) Remove(name string) {
	delete(e.Values, name)
}

func (e *Environment) DefineAbove(name string, val any) error {
	// Ensure above env exists
	if e.Enlcosing == nil {
		return fmt.Errorf("unable to define var %s to enclosing env", name)
	}
	// Define
	e.Enlcosing.Define(name, val)
	return nil
}

// Need to be able to remove from env and add to another simo
// We hoist up and down an env when we see
