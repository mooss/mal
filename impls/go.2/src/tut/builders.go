package tut

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

// Sym builds a Symbol.
func Sym(s string) typ.Symbol { return typ.Symbol(s) }

// Str builds a String.
func Str(s string) typ.String { return typ.String(s) }

// Num builds an Integer.
func Num(n int64) typ.Integer { return typ.Integer(n) }

// Nil builds Nil.
func Nil() typ.Nil { return typ.Nil{} }

// Lst builds a Lst form.
var Lst = typ.NewList

// Vec builds a VectorForm.
func Vec(elems ...core.Value) typ.VectorForm { return typ.VectorForm(elems) }

// Quote builds a quoted value, i.e. a list with two elements, the symbol quote and the given value.
func Quote(v core.Value) typ.List { return Lst(Sym("quote"), v) }

// MapForm builds a MapForm, panicking on odd length.
func MapForm(pairs ...core.Value) typ.MapForm {
	f, err := typ.NewMapForm(pairs...)
	if err != nil {
		panic(err)
	}
	return f
}

// Values gather several values into a slice.
func Values(args ...core.Value) []core.Value { return args }
