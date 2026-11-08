package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func defBang(envr core.Env, key typ.Symbol, value core.Value) core.Value {
	evald := value.Eval(envr)
	if ut.Is[typ.Fatal](evald) {
		return evald
	}

	envr.Set(string(key), evald)
	return evald
}

var DefBang = typ.SpII(defBang)

func letStar(envr core.Env, bindings typ.List, form core.Value) core.Value {
	if len(bindings)%2 != 0 {
		return typ.Fatalf("odd-numbered let bindings")
	}

	letenv := envr.Nest(nil)
	for i := 0; i < len(bindings); i += 2 {
		key, err := typ.As[typ.Symbol](bindings[i])
		if err.NotNil() {
			return err.Qualify("invalid key in let bindings")
		}

		value := bindings[i+1].Eval(letenv) //nolint:gosec
		if ut.Is[typ.Fatal](value) {
			return value
		}

		letenv.Set(string(key), value)
	}

	return form.Eval(letenv)
}

var LetStar = typ.SpII(letStar)

var MalSpecial core.Env = core.Bindings{
	"def!": DefBang,
	"let*": LetStar,
}.AsEnv()
