package envs

import (
	"fmt"
	"slices"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

var MalSpecial core.Env = core.Bindings{
	"def!": typ.SpII(defBang),
	"fn*":  typ.SpII(lambda),
	"let*": typ.SpII(letStar),
	"do":   typ.Special(specDo),
	"if":   typ.Special(specIf),
}.AsEnv()

//////////////////////////////////
// Special forms implementation //

// defBang introduces a new binding in an environment.
func defBang(envr core.Env, key typ.Symbol, value core.Value) core.Value {
	evald := value.Eval(envr)
	if ut.Is[typ.Fatal](evald) {
		return evald
	}

	envr.Set(string(key), evald)
	return evald
}

func lambda(envr core.Env, paramlist core.Value, body core.Value) core.Value {
	bindlist, err := asListOrVec(paramlist) // Assert params is a list or vector.
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
	bindlist, err := asListOrVec(bindings)
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

		value := bindlist[i+1].Eval(letenv)
		if ut.Is[typ.Fatal](value) {
			return value
		}

		letenv.Set(string(key), value)
	}

	return form.Eval(letenv)
}

func specDo(envr core.Env, args ...core.Value) core.Value {
	if len(args) == 0 {
		return typ.Nil{}
	}

	var res core.Value
	for _, arg := range args {
		res = arg.Eval(envr)
		if ut.Is[typ.Fatal](res) {
			return res
		}
	}
	return res
}

func specIfElse(envr core.Env, cond core.Value, cons core.Value, alt core.Value) core.Value {
	condeval := cond.Eval(envr)
	realised := true // Whether the condition was realised.
	switch typed := condeval.(type) {
	case typ.Fatal:
		return condeval
	case typ.Nil:
		realised = false
	case typ.Bool:
		realised = bool(typed)
	default:
	}

	if realised {
		return cons.Eval(envr)
	}

	return alt.Eval(envr)
}

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

/////////////////
// Local utils //

func asListOrVec(elements core.Value) ([]core.Value, typ.Fatal) {
	var bindlist []core.Value
	switch typed := elements.(type) {
	case typ.List:
		bindlist = typed
	case typ.VectorForm:
		bindlist = typed
	default:
		return nil, typ.Expected(fmt.Sprintf("%s or %s", core.LIST, core.VECTOR_FORM), elements)
	}
	return bindlist, typ.Fatal{}
}

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
		return body.Eval(envr.Nest(bindings))
	}
}

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

		return body.Eval(envr.Nest(bindings))
	}
}
