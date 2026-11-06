package core

import (
	"fmt"
	"strconv"
	"strings"
)

type TypeName string

//revive:disable:var-naming
const (
	SYMBOL      TypeName = "symbol"
	LIST        TypeName = "list"
	FUNCTION    TypeName = "function"
	INTEGER     TypeName = "integer"
	STRING      TypeName = "string"
	VECTOR      TypeName = "vector"
	MAP         TypeName = "map"
	MAP_LITERAL TypeName = "map-literal"
	ERROR       TypeName = "error"
	NIL         TypeName = "nil"
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
type Map map[Type]Type
type MapLit []Type
type Error struct{ error }
type Nil struct{}

func (Symbol) Type() TypeName   { return SYMBOL }
func (List) Type() TypeName     { return LIST }
func (Function) Type() TypeName { return FUNCTION }
func (Integer) Type() TypeName  { return INTEGER }
func (String) Type() TypeName   { return STRING }
func (Vector) Type() TypeName   { return VECTOR }
func (Map) Type() TypeName      { return MAP }
func (MapLit) Type() TypeName   { return MAP_LITERAL }
func (Error) Type() TypeName    { return ERROR }
func (Nil) Type() TypeName      { return NIL }

func (s Symbol) String() string  { return string(s) }
func (l List) String() string    { return "(" + strings.Join(_map(Type.String, l), " ") + ")" }
func (Function) String() string  { return "<function>" }
func (i Integer) String() string { return strconv.FormatInt(int64(i), 10) }
func (s String) String() string  { return strconv.Quote(string(s)) }
func (v Vector) String() string  { return "[" + strings.Join(_map(Type.String, v), " ") + "]" }
func (m Map) String() string     { return "{" + strings.Join(mapflat(Type.String, m), " ") + "}" }
func (m MapLit) String() string  { return "{" + strings.Join(_map(Type.String, m), " ") + "}" }
func (e Error) String() string   { return fmt.Sprintf("(error %q)", e.Error()) }
func (Nil) String() string       { return "nil" }

func _map[From, To any](fun func(From) To, source []From) []To {
	res := make([]To, len(source))
	for i, el := range source {
		res[i] = fun(el)
	}
	return res
}

func mapflat[From comparable, To any](fun func(From) To, source map[From]From) []To {
	res := make([]To, len(source)*2)
	i := 0
	for k, v := range source {
		res[i] = fun(k)
		i++
		res[i] = fun(v)
		i++
	}
	return res
}

func Errorf(msg string, args ...any) Error {
	return Error{error: fmt.Errorf(msg, args...)}
}
