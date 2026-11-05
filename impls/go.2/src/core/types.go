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
	VECTOR   TypeName = "vector"
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
type Vector []Type
type Error struct{ error }
type Dummy struct{}

func (Symbol) Type() TypeName   { return SYMBOL }
func (List) Type() TypeName     { return LIST }
func (Function) Type() TypeName { return FUNCTION }
func (Integer) Type() TypeName  { return INTEGER }
func (Vector) Type() TypeName   { return VECTOR }
func (Error) Type() TypeName    { return ERROR }
func (Dummy) Type() TypeName    { return DUMMY }

func (s Symbol) String() string  { return string(s) }
func (l List) String() string    { return "(" + strings.Join(_map(Type.String, l), " ") + ")" }
func (Function) String() string  { return "<function>" }
func (i Integer) String() string { return strconv.FormatInt(int64(i), 10) }
func (v Vector) String() string  { return "[" + strings.Join(_map(Type.String, v), " ") + "]" }
func (e Error) String() string   { return fmt.Sprintf("(error %q)", e.Error()) }
func (Dummy) String() string     { return "dummy" }

func _map[From, To any](fun func(From) To, source []From) []To {
	res := make([]To, len(source))
	for i, el := range source {
		res[i] = fun(el)
	}
	return res
}

func Errorf(msg string, args ...any) Error {
	return Error{error: fmt.Errorf(msg, args...)}
}
