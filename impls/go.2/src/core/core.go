// Package core defines fundamental types for the interpreter.
package core

import "maps"

// TypeName identifies a Mal type.
type TypeName string

//revive:disable:var-naming
const (
	// FATAL is the name of the fatal error type.
	FATAL TypeName = "fatal"
	// FUNCTION is the name of the function type.
	FUNCTION TypeName = "function"
	// INTEGER is the name of the integer type.
	INTEGER TypeName = "integer"
	// LIST is the name of the list type.
	LIST TypeName = "list"
	// MAP is the name of the map type.
	MAP TypeName = "map"
	// MAP_FORM is the name of the unevaluated map type.
	MAP_FORM TypeName = "map-form"
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

// Bindings is a primitive environment associating symbols to their values.
type Bindings map[string]Value

func (b Bindings) AsEnv() Env { return Env{inner: b} }

// Env is an environment associating symbols to their values.
type Env struct {
	inner Bindings
}

// NewEnv creates a new environment, optionally pre-populated with the given bindings.
func NewEnv(bindings ...Bindings) Env {
	store := make(map[string]Value)
	for _, b := range bindings {
		maps.Copy(store, b)
	}
	return Env{inner: store}
}

// Get returns the value bound to key, reporting whether it was present.
func (e Env) Get(key string) (Value, bool) {
	value, present := e.inner[key]
	return value, present
}

// Set binds key to val in the environment.
func (e Env) Set(key string, val Value) {
	if e.inner == nil {
		e.inner = make(map[string]Value)
	}
	e.inner[key] = val
}
