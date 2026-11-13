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
func Lst(elems ...core.Value) typ.List { return typ.List(elems) }

// Vec builds a VectorForm.
func Vec(elems ...core.Value) typ.VectorForm { return typ.VectorForm(elems) }

// Quote builds a Quote.
func Quote(v core.Value) typ.Quote { return typ.Quote{Value: v} }

// MapForm builds a MapForm, panicking on odd length.
func MapForm(pairs ...core.Value) typ.MapForm {
	f, err := typ.NewMapForm(pairs...)
	if err != nil {
		panic(err)
	}
	return f
}
