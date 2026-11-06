// Package core implements lisp types.
package core

// TypeName identifies a Mal type.
type TypeName string

//revive:disable:var-naming
const (
	// SYMBOL is the name of the symbol type.
	SYMBOL TypeName = "symbol"
	// LIST is the name of the list type.
	LIST TypeName = "list"
	// FUNCTION is the name of the function type.
	FUNCTION TypeName = "function"
	// INTEGER is the name of the integer type.
	INTEGER TypeName = "integer"
	// STRING is the name of the string type.
	STRING TypeName = "string"
	// VECTOR is the name of the vector type.
	VECTOR TypeName = "vector"
	// MAP is the name of the map type.
	MAP TypeName = "map"
	// MAP_LITERAL is the name of the unevaluated map type.
	MAP_LITERAL TypeName = "map-literal"
	// ERROR is the name of the error type.
	ERROR TypeName = "error"
	// NIL is the name of the nil type.
	NIL TypeName = "nil"
)

// Type is a Mal value.
type Type interface {
	Type() TypeName
	String() string
}
