package tut

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

// Kyw builds a keyword.
func Kyw(s string) typ.Keyword { return typ.Keyword(s) }

// Lst builds a Lst form.
var Lst = typ.NewList

// MapForm builds a MapForm, panicking on odd length.
func MapForm(pairs ...core.Value) typ.MapForm {
	f, err := typ.NewMapForm(pairs...)
	if err != nil {
		panic(err)
	}
	return f
}

// Nil builds Nil.
func Nil() typ.Nil { return typ.Nil{} }

// Int builds an Integer.
func Int(n int64) typ.Integer { return typ.Integer(n) }

// Str builds a String.
func Str(s string) typ.String { return typ.NewString(s) }

// Sym builds a Symbol.
func Sym(s string) typ.Symbol { return typ.Symbol(s) }

// Values gather several values into a slice.
func Values(args ...core.Value) []core.Value { return args }

// Vec builds a VectorForm.
func Vec(elems ...core.Value) typ.VectorForm { return typ.VectorForm(elems) }
