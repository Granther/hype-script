package glorpups

import (
	"fmt"
	"hype-script/internal/token"
)

type Glorpup interface {
	Error() string
}

type RuntimeGlorpup struct {
	Token    token.Token
	Message  string
	Previous Glorpup
}

type IndexBoundsGlorpup struct {
	Token    token.Token
	Message  string
	Previous Glorpup
}

type TypeGlorpup struct {
	Token    token.Token
	Message  string
	Previous Glorpup
}

func NewRuntimeGlorpup(token token.Token, message string, err Glorpup) Glorpup {
	return &RuntimeGlorpup{
		Token:    token,
		Message:  message,
		Previous: err,
	}
}

func (g *RuntimeGlorpup) Error() string {
	return Report(g.Message, g.Previous)
}

func NewTypeGlorpup(token token.Token, message string, err Glorpup) Glorpup {
	return &TypeGlorpup{
		Token:    token,
		Message:  message,
		Previous: err,
	}
}

func (g *TypeGlorpup) Error() string {
	return Report(g.Message, g.Previous)
}

func NewIndexBoundsGlorpup(token token.Token, message string, err Glorpup) Glorpup {
	return &IndexBoundsGlorpup{
		Token: token,
		Message: message,
		Previous: err,
	}
}

func (g *IndexBoundsGlorpup) Error() string {
	return Report(g.Message, g.Previous)
}

func InterpreterRuntimeError(message string, err Glorpup) {
	fmt.Println(Report(message, err))
}

func Report(message string, err Glorpup) string {
	if err == nil {
		return message
	}
	return fmt.Sprintf("%s\n ^\n |\n%s\n", message, err.Error())
}
