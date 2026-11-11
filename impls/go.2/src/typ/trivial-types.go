// This file implements trivial types, i.e. types whose representation and evaluation amounts to
// basically nothing.

package typ

import (
	"fmt"
	"strconv"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/ut"
)

// Bool is a boolean.
type Bool bool

var True = Bool(true)
var False = Bool(false)

func (b Bool) Equal(other core.Value) bool { return TrivialEq(b, other) }
func (b Bool) Eval(core.Env) core.Value    { return b }
func (b Bool) String() string              { return strconv.FormatBool(bool(b)) }
func (Bool) Type() core.TypeName           { return core.BOOL }

// Fatal is an error that interrupts all execution immediately and irrecoverably.
type Fatal struct{ error }

// Fatalf creates a Fatal error with a formatted message.
func Fatalf(msg string, args ...any) Fatal { return Fatal{error: fmt.Errorf(msg, args...)} }

// Equal returns true if both are nil errors.
func (f Fatal) Equal(other core.Value) bool {
	o, ok := other.(Fatal)
	return ok && f.error == nil && o.error == nil
}

func (f Fatal) Eval(core.Env) core.Value { return f }

// NotNil returns true when the embedded error is not nil.
func (f Fatal) NotNil() bool { return f.error != nil }

// Qualify returns a new fatal error with "{prefix}: " added in front of the error.
func (f Fatal) Qualify(prefix string) Fatal {
	return Fatal{error: fmt.Errorf("%s: %w", prefix, f.error)}
}

func (f Fatal) String() string    { return fmt.Sprintf("(fatal %q)", f.Error()) }
func (Fatal) Type() core.TypeName { return core.FATAL }

// Integer is a 64-bit signed integer value.
type Integer int64

func (i Integer) Equal(other core.Value) bool { return TrivialEq(i, other) }
func (i Integer) Eval(core.Env) core.Value    { return i }
func (i Integer) String() string              { return strconv.FormatInt(int64(i), 10) }
func (Integer) Type() core.TypeName           { return core.INTEGER }

// Nil is the Mal nil value.
type Nil struct{}

func (Nil) Equal(other core.Value) bool { return ut.Is[Nil](other) }
func (n Nil) Eval(core.Env) core.Value  { return n }
func (Nil) String() string              { return "nil" }
func (Nil) Type() core.TypeName         { return core.NIL }

// Quote is a quoted value which returns the raw value upon evaluation.
type Quote struct{ Value core.Value }

func (q Quote) Equal(other core.Value) bool {
	o, ok := other.(Quote)
	return ok && q.Value.Equal(o.Value)
}

func (q Quote) Eval(core.Env) core.Value { return q.Value }
func (q Quote) String() string           { return "(quote " + q.Value.String() + ")" }
func (Quote) Type() core.TypeName        { return core.QUOTE }

// String is a Mal string value.
type String string

func (s String) Equal(other core.Value) bool { return TrivialEq(s, other) }
func (s String) Eval(core.Env) core.Value    { return s }
func (s String) String() string              { return strconv.Quote(string(s)) }
func (String) Type() core.TypeName           { return core.STRING }

// Symbol is a Mal symbol value.
type Symbol string

func (s Symbol) Equal(other core.Value) bool { return TrivialEq(s, other) }

func (s Symbol) Eval(envr core.Env) core.Value {
	value := envr.Get(string(s))
	if value == nil {
		return Fatalf("symbol '%s' not found", s)
	}
	return value
}

func (s Symbol) String() string    { return string(s) }
func (Symbol) Type() core.TypeName { return core.SYMBOL }

// Keyword is a Mal keyword value, written with a leading colon.
type Keyword string

func (k Keyword) Equal(other core.Value) bool { return TrivialEq(k, other) }
func (k Keyword) Eval(core.Env) core.Value    { return k }
func (k Keyword) String() string              { return ":" + string(k) }
func (Keyword) Type() core.TypeName           { return core.TypeName("keyword") }
