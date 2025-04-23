package native

import (
	"time"
	"glorp/types"
)

type ClockCallable struct{}

func NewClockCallable() Callable {
	return &ClockCallable{}
}

func (c *ClockCallable) Call(interpreter types.Interpreter, args []any) (any, error) {
	return time.Now().UnixNano() / 1e9, nil
}

func (c *ClockCallable) Arity() int {
	return 0
}

func (c *ClockCallable) String() string {
	return "<native fn>"
}
