// Package core defines fundamental types for the interpreter.
package core

// TypeName identifies a Mal type.
type TypeName string

//revive:disable:var-naming
const (
	// ATOM tags atom, a.k.a. references.
	ATOM TypeName = "atom"
	// BOOL tags boolean value.
	BOOL TypeName = "bool"
	// FATAL tags fatal error aborting evaluation.
	FATAL TypeName = "fatal"
	// FUNCTION tags callable function value.
	FUNCTION TypeName = "function"
	// INTEGER tags integer number.
	INTEGER TypeName = "integer"
	// LIST tags list of values.
	LIST TypeName = "list"
	// MAP tags map of keys to values.
	MAP TypeName = "map"
	// MAP_FORM tags unevaluated map awaiting evaluation.
	MAP_FORM TypeName = "map-form"
	// NIL tags nil value representing absence.
	NIL TypeName = "nil"
	// QUOTE tags quoted form preventing evaluation.
	QUOTE TypeName = "quote"
	// SPECIAL tags special form evaluator.
	SPECIAL TypeName = "special"
	// STRING tags string of characters.
	STRING TypeName = "string"
	// SYMBOL tags symbolic identifier.
	SYMBOL TypeName = "symbol"
	// TAILCALL tags deferred computation for tail call optimization.
	TAILCALL TypeName = "tail-call"
	// VECTOR tags vector of values.
	VECTOR TypeName = "vector"
	// VECTOR_FORM tags unevaluated vector awaiting evaluation.
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
	// Equal returns true when other is equal to the receiver.
	Equal(Value) bool
}
