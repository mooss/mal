// This file implements trivial types, i.e. types whose representation and evaluation amounts to
// basically nothing.

package typ

import (
	"fmt"
	"strconv"

	"github.com/mooss/mal/src/core"
)

// Fatal is an error that interrupts all execution immediately and irrecoverably.
type Fatal struct{ error }

// Fatalf creates a Fatal error with a formatted message.
func Fatalf(msg string, args ...any) Fatal { return Fatal{error: fmt.Errorf(msg, args...)} }
func (Fatal) Type() core.TypeName          { return core.ERROR }
func (e Fatal) String() string             { return fmt.Sprintf("(fatal %q)", e.Error()) }
func (e Fatal) Eval(core.Env) core.Value   { return e }

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
	value, present := envr[string(s)]
	if !present {
		return Fatalf("symbol %#q not found", s)
	}
	return value
}
