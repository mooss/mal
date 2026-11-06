// Package core defines the basic lisp Type interface.
package core

// TypeName identifies a Mal type.
type TypeName string

//revive:disable:var-naming
const (
	// ERROR is the name of the error type.
	ERROR TypeName = "error"
	// FUNCTION is the name of the function type.
	FUNCTION TypeName = "function"
	// INTEGER is the name of the integer type.
	INTEGER TypeName = "integer"
	// LIST is the name of the list type.
	LIST TypeName = "list"
	// MAP is the name of the map type.
	MAP TypeName = "map"
	// MAP_LITERAL is the name of the unevaluated map type.
	MAP_LITERAL TypeName = "map-literal"
	// NIL is the name of the nil type.
	NIL TypeName = "nil"
	// STRING is the name of the string type.
	STRING TypeName = "string"
	// SYMBOL is the name of the symbol type.
	SYMBOL TypeName = "symbol"
	// VECTOR is the name of the vector type.
	VECTOR TypeName = "vector"
	// VECTOR_FORM is the name of the unevaluated vector type.
	VECTOR_FORM TypeName = "vector-form"
)

// Value defines the methods a Mal value must implement.
type Value interface {
	// Type returns the TypeName identifying the type of the value.
	Type() TypeName
	// String returns the string representation of the Mal value as Mal code.
	String() string
	// Eval evaluates the value against an environment.
	Eval(Env) Value
}

// Env is an environment associating symbols to their values.
type Env map[string]Value
