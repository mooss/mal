// Package core implements lisp types.
package core

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
