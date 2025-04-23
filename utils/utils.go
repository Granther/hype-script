package utils

import (
	"fmt"
	"glorp/types"
	"strings"
)

// If a and b can be converted to ints, return converted and true
func ConvFloat(a, b any) (float64, float64, bool) {
	left, lok := a.(float64)
	right, rok := b.(float64)

	if lok && rok {
		return left, right, true
	} else {
		return 0, 0, false
	}
}

func Parenthesize(visitor types.Visitor, name string, exprs ...types.Expr) (string, error) {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		val, err := expr.Accept(visitor)
		if err != nil {
			return "", err
		}
		builder.WriteString(fmt.Sprintf("%v", val))
	}
	builder.WriteString(")")

	return builder.String(), nil
}

func IsFloat(val any) (float64, bool) {
	ival, ok := val.(float64)
	if !ok {
		return 0, false
	}
	return ival, true
}

func Stringify(val any) string {
	if val == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", val)
}
