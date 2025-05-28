package native

import (
	"time"
	"hype-script/internal/types/core"
)

type ClockCallable struct{}

func NewClockCallable() Callable {
	return &ClockCallable{}
}

func (c *ClockCallable) Call(interpreter core.InterpreterHandler, args []any) (any, error) {
	return time.Now().UnixNano() / 1e9, nil
}

func (c *ClockCallable) Arity() int {
	return 0
}

func (c *ClockCallable) String() string {
	return "<native fn>"
}
