package literal

import "fmt"

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

