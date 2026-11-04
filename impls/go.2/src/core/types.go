package core

import (
	"fmt"
	"strconv"
	"strings"
)

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
	String() string
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

func (s Symbol) String() string  { return string(s) }
func (i Integer) String() string { return strconv.FormatInt(int64(i), 10) }
func (l List) String() string {
	return "(" + strings.Join(l.Map(func(t Type) string { return t.String() }), " ") + ")"
}
func (e Error) String() string  { return fmt.Sprintf("(error %q)", e.Error()) }
func (Function) String() string { return "<function>" }
func (Dummy) String() string    { return "dummy" }

func (l List) Map[To any](fun func(Type) To) []To {
	res := make([]To, len(l))
	for i, el := range l {
		res[i] = fun(el)
	}
	return res
}

func Errorf(msg string, args ...any) Error {
	return Error{error: fmt.Errorf(msg, args...)}
}
