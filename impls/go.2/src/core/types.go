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
	STRING   TypeName = "string"
	VECTOR   TypeName = "vector"
	ERROR    TypeName = "error"
	DUMMY    TypeName = "dummy"
	NIL      TypeName = "nil"
)

type Type interface {
	Type() TypeName
	String() string
}

type Symbol string
type List []Type
type Function func(Integer, Integer) Integer
type Integer int64
type String string
type Vector []Type
type Error struct{ error }
type Dummy struct{}
type Nil struct{}

func (Symbol) Type() TypeName   { return SYMBOL }
func (List) Type() TypeName     { return LIST }
func (Function) Type() TypeName { return FUNCTION }
func (Integer) Type() TypeName  { return INTEGER }
func (String) Type() TypeName   { return STRING }
func (Vector) Type() TypeName   { return VECTOR }
func (Error) Type() TypeName    { return ERROR }
func (Dummy) Type() TypeName    { return DUMMY }
func (Nil) Type() TypeName      { return NIL }

func (s Symbol) String() string  { return string(s) }
func (l List) String() string    { return "(" + strings.Join(_map(Type.String, l), " ") + ")" }
func (Function) String() string  { return "<function>" }
func (i Integer) String() string { return strconv.FormatInt(int64(i), 10) }
func (s String) String() string  { return strconv.Quote(string(s)) }
func (v Vector) String() string  { return "[" + strings.Join(_map(Type.String, v), " ") + "]" }
func (e Error) String() string   { return fmt.Sprintf("(error %q)", e.Error()) }
func (Dummy) String() string     { return "dummy" }
func (Nil) String() string       { return "nil" }

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
