// This file implements trivial types, i.e. types whose representation and evaluation amounts to
// basically nothing.

package typ

import (
	"fmt"
	"strconv"

	"github.com/mooss/mal/src/core"
)

// Bool is a boolean.
type Bool bool

func (Bool) Type() core.TypeName        { return core.BOOL }
func (b Bool) String() string           { return strconv.FormatBool(bool(b)) }
func (b Bool) Eval(core.Env) core.Value { return b }

// Fatal is an error that interrupts all execution immediately and irrecoverably.
type Fatal struct{ error }

// Fatalf creates a Fatal error with a formatted message.
func Fatalf(msg string, args ...any) Fatal { return Fatal{error: fmt.Errorf(msg, args...)} }
func (Fatal) Type() core.TypeName          { return core.FATAL }
func (f Fatal) String() string             { return fmt.Sprintf("(fatal %q)", f.Error()) }
func (f Fatal) Eval(core.Env) core.Value   { return f }

// NotNil returns true when the embedded error is not nil.
func (f Fatal) NotNil() bool { return f.error != nil }

// Qualify returns a new fatal error with "{prefix}: " added in front of the error.
func (f Fatal) Qualify(prefix string) Fatal {
	return Fatal{error: fmt.Errorf("%s: %w", prefix, f.error)}
}

// Integer is a 64-bit signed integer value.
type Integer int64

func (Integer) Type() core.TypeName        { return core.INTEGER }
func (i Integer) String() string           { return strconv.FormatInt(int64(i), 10) }
func (i Integer) Eval(core.Env) core.Value { return i }

// Nil is the Mal nil value.
type Nil struct{}

func (Nil) Type() core.TypeName        { return core.NIL }
func (Nil) String() string             { return "nil" }
func (n Nil) Eval(core.Env) core.Value { return n }

// Quote is a quoted value which returns the raw value upon evaluation.
type Quote struct{ Value core.Value }

func (Quote) Type() core.TypeName        { return core.QUOTE }
func (q Quote) String() string           { return "'" + q.Value.String() }
func (q Quote) Eval(core.Env) core.Value { return q.Value }

// String is a Mal string value.
type String string

func (String) Type() core.TypeName        { return core.STRING }
func (s String) String() string           { return strconv.Quote(string(s)) }
func (s String) Eval(core.Env) core.Value { return s }

// Symbol is a Mal symbol value.
type Symbol string

func (Symbol) Type() core.TypeName { return core.SYMBOL }
func (s Symbol) String() string    { return string(s) }
func (s Symbol) Eval(envr core.Env) core.Value {
	value := envr.Get(string(s))
	if value == nil {
		return Fatalf("symbol '%s' not found", s)
	}
	return value
}
