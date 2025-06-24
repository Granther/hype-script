package literal

import "fmt"

// Simply a container for a value
// On a more complex note, a way to encode the value of a written in value from the src to a var in the backend

type Literal struct {
	Val any
}

func NewLiteral(val any) *Literal {
	return &Literal{
		Val: val,
	}
}

func (l *Literal) String() string {
	return fmt.Sprintf("%v", l.Val)
}
