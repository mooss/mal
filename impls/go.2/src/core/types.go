package core

import "fmt"

type TypeName string

const (
	SYMBOL   TypeName = "symbol"
	LIST     TypeName = "list"
	FUNCTION TypeName = "function"
	INTEGER  TypeName = "integer"
	ERROR    TypeName = "error"
	DUMMY    TypeName = "dummy"
)

type Type interface {
	Type() TypeName
}

type Symbol string
type List []Type
type Function func(Integer, Integer) Integer
type Integer int64
type Error struct{ error }
type Dummy struct{}

func (Symbol) Type() TypeName   { return SYMBOL }
func (List) Type() TypeName     { return LIST }
func (Function) Type() TypeName { return FUNCTION }
func (Integer) Type() TypeName  { return INTEGER }
func (Error) Type() TypeName    { return ERROR }
func (Dummy) Type() TypeName    { return DUMMY }

func (lst List) Map[To any](fun func(Type) To) []To {
	res := make([]To, len(lst))
	for i, el := range lst {
		res[i] = fun(el)
	}
	return res
}

func Errorf(msg string, args ...any) Error {
	return Error{error: fmt.Errorf(msg, args...)}
}
