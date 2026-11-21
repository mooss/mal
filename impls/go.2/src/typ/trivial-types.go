// This file implements trivial types, i.e. types whose representation and evaluation amounts to
// basically nothing.

package typ

import (
	"strconv"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/strfun"
	"github.com/mooss/mal/src/ut"
)

// An atom is a reference to a value.
// It can be manipulated with functions like `atom`, `atom?`, `deref`, `reset!` and `swap!`.
type Atom struct{ Ref core.Value }

func (a *Atom) Equal(other core.Value) bool { return TrivialEq(a, other) }
func (a *Atom) Eval(core.Env) core.Value    { return a }
func (a *Atom) String() string              { return "(atom " + a.Ref.String() + ")" }
func (*Atom) Type() core.TypeName           { return core.ATOM }

// Bool is a boolean.
type Bool bool

var True = Bool(true)
var False = Bool(false)

func (b Bool) Equal(other core.Value) bool { return TrivialEq(b, other) }
func (b Bool) Eval(core.Env) core.Value    { return b }
func (b Bool) String() string              { return strconv.FormatBool(bool(b)) }
func (Bool) Type() core.TypeName           { return core.BOOL }

// Integer is a 64-bit signed integer value.
type Integer int64

func (i Integer) Equal(other core.Value) bool { return TrivialEq(i, other) }
func (i Integer) Eval(core.Env) core.Value    { return i }
func (i Integer) String() string              { return strconv.FormatInt(int64(i), 10) }
func (Integer) Type() core.TypeName           { return core.INTEGER }
func (i Integer) Unbox() int64                { return int64(i) }

// Keyword is a Mal keyword value, written with a leading colon.
type Keyword string

func (k Keyword) Equal(other core.Value) bool { return TrivialEq(k, other) }
func (k Keyword) Eval(core.Env) core.Value    { return k }
func (k Keyword) String() string              { return ":" + string(k) }
func (Keyword) Type() core.TypeName           { return core.TypeName("keyword") }

// Nil is the Mal nil value.
type Nil struct{}

func (Nil) Equal(other core.Value) bool { return ut.Is[Nil](other) }
func (n Nil) Eval(core.Env) core.Value  { return n }
func (Nil) String() string              { return "nil" }
func (Nil) Type() core.TypeName         { return core.NIL }

// String is a Mal string value.
type String string

func (s String) Equal(other core.Value) bool { return TrivialEq(s, other) }
func (s String) Eval(core.Env) core.Value    { return s }
func (s String) String() string              { return strfun.MalQuote(string(s)) }
func (String) Type() core.TypeName           { return core.STRING }

// Symbol is a Mal symbol value.
type Symbol string

func (s Symbol) Equal(other core.Value) bool { return TrivialEq(s, other) }

func (s Symbol) Eval(envr core.Env) core.Value {
	value := envr.Get(string(s))
	if value == nil {
		return Throwf("'%s' not found", s)
	}
	return value
}

func (s Symbol) String() string    { return string(s) }
func (Symbol) Type() core.TypeName { return core.SYMBOL }
