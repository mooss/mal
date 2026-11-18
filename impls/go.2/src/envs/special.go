package envs

import (
	"slices"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/lazy"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

var MalSpecial core.Env = core.Bindings{
	"def!":       typ.SpII(defBang),
	"do":         typ.Special(specDo),
	"fn*":        typ.SpII(lambda),
	"if":         typ.Special(specIf),
	"let*":       typ.SpII(letStar),
	"quasiquote": typ.SpI(quasiquote),
	"quote":      typ.SpI(quote),
}.AsEnv()

// ExtraSpecial contains special forms that are not a part of Mal.
var ExtraSpecial core.Env = core.Bindings{
	"env": typ.SpZ(env),
}.AsEnv()

//////////////////////////////////
// Special forms implementation //

// defBang introduces a new binding in an environment.
func defBang(envr core.Env, key typ.Symbol, value core.Value) core.Value {
	evald := typ.Resolve(envr, value)
	if ut.Is[typ.Fatal](evald) {
		return evald
	}

	envr.Set(string(key), evald)
	return evald
}

// lambda builds a closure capturing the current environment, with fixed or variadic parameters.
// A `&` symbol before the last parameter marks it as capturing the remaining arguments.
func lambda(envr core.Env, paramlist core.Value, body core.Value) core.Value {
	bindlist, err := typ.AsListOrVec(paramlist) // Assert params is a list or vector.
	if err.NotNil() {
		return err
	}

	symlist, err := typ.AsSlice[typ.Symbol](bindlist...) // Assert all elements are symbol.
	if err.NotNil() {
		return err
	}

	amperpos := slices.Index(symlist, typ.Symbol("&"))
	if amperpos < 0 {
		return typ.Function(mkFixedClosure(envr, symlist, body))
	}
	if amperpos != len(symlist)-2 {
		return typ.Fatalf("invalid & in lambda params, must be right before the last symbol")
	}

	rest := string(symlist[amperpos+1])
	return typ.Function(mkVariadicClosure(envr, symlist[:amperpos], rest, body))
}

// letStar creates a temporary environment containing the given bindings and evaluates form in this
// new environment.
func letStar(envr core.Env, bindings core.Value, form core.Value) core.Value {
	bindlist, err := typ.AsListOrVec(bindings)
	if err.NotNil() {
		return err
	}

	if len(bindlist)%2 != 0 {
		return typ.Fatalf("odd-numbered let bindings")
	}

	letenv := envr.Nest(nil)
	for i := 0; i < len(bindlist); i += 2 {
		key, err := typ.As[typ.Symbol](bindlist[i])
		if err.NotNil() {
			return err.Qualify("invalid key in let bindings")
		}

		value := typ.Resolve(letenv, bindlist[i+1])
		if ut.Is[typ.Fatal](value) {
			return value
		}

		letenv.Set(string(key), value)
	}

	return typ.Tail(letenv, form)
}

// MkEval returns an eval function closed over envr.
func MkEval(envr core.Env) func(core.Value) core.Value {
	return func(ast core.Value) core.Value {
		return typ.Resolve(envr, ast)
	}
}

func quasiquote(envr core.Env, ast core.Value) core.Value {
	return typ.Tail(envr, quasiquoteImpl(ast))
}

func quasiquoteImpl(ast core.Value) core.Value {
	switch typed := ast.(type) {
	case typ.List:
		return quasiquoteList(typed)
	case typ.Vector:
		return typ.List{typ.Symbol("vec"), quasiquotePlainList(typ.List(typed))}
	case typ.Map, typ.MapForm, typ.Symbol:
		return typ.Quote(typed)
	default:
		return ast
	}
}

//revive:disable:cognitive-complexity
func quasiquoteList(ast typ.List) core.Value {
	if len(ast) == 0 {
		return ast
	}

	if typ.Symbol("unquote").Equal(ast[0]) {
		if err := typ.ExpectNArgs(2, ast); err.NotNil() {
			return err.Qualify("unquote")
		}
		return ast[1]
	}

	return quasiquotePlainList(ast)
}

// quasiquotePlainList processes through quasiquote a list that does not start with unquote.
func quasiquotePlainList(ast typ.List) core.Value {
	res := typ.List{}
	for _, elt := range slices.Backward(ast) {
		eltlst, isList := elt.(typ.List)
		if isList && len(eltlst) > 0 && typ.Symbol("splice-unquote").Equal(eltlst[0]) {
			if err := typ.ExpectNArgs(2, eltlst); err.NotNil() {
				return err.Qualify("splice-unquote")
			}
			res = typ.List{typ.Symbol("concat"), eltlst[1], res}
			continue
		}

		res = typ.List{typ.Symbol("cons"), quasiquoteImpl(elt), res}
	}
	return res
}

// quote returns an expression as-is, without any evaluation.
func quote(_ core.Env, ast core.Value) core.Value { return ast }

// specDo evaluates each argument in order and returns the value of the last one.
func specDo(envr core.Env, args ...core.Value) core.Value {
	if len(args) == 0 {
		return typ.Nil{}
	}

	for _, arg := range args[:len(args)-1] {
		res := typ.Resolve(envr, arg)
		if ut.Is[typ.Fatal](res) {
			return res
		}
	}

	return typ.Tail(envr, args[len(args)-1])
}

// specIfElse evaluates the condition and returns the evaluation of `cons` if it is truthy,
// otherwise the evaluation of `alt`.
// Only nil and false are falsy.
func specIfElse(envr core.Env, cond core.Value, cons core.Value, alt core.Value) core.Value {
	condeval := typ.Resolve(envr, cond)
	realised, err := typ.Truthful(condeval)
	if err.NotNil() {
		return err
	}

	if realised {
		return typ.Tail(envr, cons)
	}

	return typ.Tail(envr, alt)
}

// specIf is the `if` special form, accepting a condition with 1 or 2 branches.
func specIf(envr core.Env, args ...core.Value) core.Value {
	switch len(args) {
	case 2:
		return specIfElse(envr, args[0], args[1], typ.Nil{})
	case 3:
		return specIfElse(envr, args[0], args[1], args[2])
	default:
		return typ.FmtExpectedNArgs("2 or 3", len(args))
	}
}

/////////////////////////
// Extra special forms //
// Those are some utilities sprinkled on top of Mal.

func env(envr core.Env) typ.Vector {
	mkassoc := func(key string, value core.Value) core.Value {
		return typ.Vector{typ.String(key), value}
	}
	return lazy.CollectHint(
		lazy.Map2(mkassoc, envr.IterAssocs()),
		envr.NAssocs(),
	)
}

/////////////////
// Local utils //

// mkFixedClosure returns a closure binding exactly len(symlist) arguments to the given symbols.
func mkFixedClosure(
	envr core.Env, symlist []typ.Symbol, body core.Value,
) func(args ...core.Value) core.Value {
	return func(args ...core.Value) core.Value {
		if err := typ.ExpectNArgs(len(symlist), args); err.NotNil() {
			return err
		}

		bindings := core.Bindings{}
		for i, sym := range symlist {
			bindings[string(sym)] = args[i]
		}
		return typ.Tail(envr.Nest(bindings), body)
	}
}

// mkVariadicClosure returns a closure binding at least len(symlist) arguments, with the remaining
// arguments collected in a list bound to restKey.
func mkVariadicClosure(
	envr core.Env, symlist []typ.Symbol, restKey string, body core.Value,
) func(args ...core.Value) core.Value {
	return func(args ...core.Value) core.Value {
		if err := typ.ExpectAtLeastNArgs(len(symlist), args); err.NotNil() {
			return err
		}

		bindings := core.Bindings{}
		for i, sym := range symlist {
			bindings[string(sym)] = args[i]
		}

		rest := make(typ.List, len(args)-len(symlist))
		for i, val := range args[len(symlist):] {
			rest[i] = val
		}
		bindings[restKey] = rest

		return typ.Tail(envr.Nest(bindings), body)
	}
}
