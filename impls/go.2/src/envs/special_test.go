//revive:disable:cognitive-complexity,cyclomatic,function-length
//nolint:dupword
package envs_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
	"github.com/mooss/mal/src/typ"
)

func TestLetStarSequentialScopingAndOuterIsolation(t *testing.T) {
	for _, input := range []string{"(let* (x a y x) y)", "(let* [x a y x] y)"} {
		outer := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("outer-a-val")})
		assert.Eval(t, outer, input, tut.Sym("outer-a-val"))
		assert.EnvLacks(t, outer, "x", "y")
		assert.EnvHas(t, outer, "a", tut.Sym("outer-a-val"))
	}
}

func TestLetStarOddBindingsError(t *testing.T) {
	for _, input := range []string{"(let* (x) x)", "(let* [x] x)"} {
		_ = assert.EvalFatalContains(t, envs.MalSpecial.Nest(nil), input, "odd")
	}
}

func TestLetStarNonSymbolKeyQualifyPrefix(t *testing.T) {
	for _, input := range []string{"(let* (() a) a)", "(let* [() a] a)"} {
		outer := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("a-val")})
		_ = assert.EvalFatalContains(t, outer, input, "invalid key in let bindings")
	}
}

func TestLetStarFatalValueForwardingWithoutSet(t *testing.T) {
	for _, input := range []string{"(let* (ok missing) ok)", "(let* [ok missing] ok)"} {
		outer := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
		assert.EvalFatal(t, outer, input)
		assert.EnvLacks(t, outer, "ok")
		assert.EnvHas(t, outer, "good", tut.Sym("good-val"))
	}
}

func TestDefBangEvalThenSetAndReturn(t *testing.T) {
	env := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("a-val")})
	assert.Eval(t, env, "(def! k a)", tut.Sym("a-val"))
	assert.EnvHas(t, env, "k", tut.Sym("a-val"))
}

func TestDefBangFatalNonSet(t *testing.T) {
	env := envs.MalSpecial.Nest(nil)
	assert.EvalFatal(t, env, "(def! k missing)")
	assert.EnvLacks(t, env, "k")
}

func TestFnStarIdentityListVecEmptyMulti(t *testing.T) {
	cases := []struct {
		input string
		want  core.Value
	}{
		{"((fn* (x) x) a)", tut.Sym("a-val")},
		{"((fn* [x] x) a)", tut.Sym("a-val")},
		{"((fn* () a))", tut.Sym("a-val")},
		{"((fn* [] a))", tut.Sym("a-val")},
		{"((fn* (x y) x) a b)", tut.Sym("a-val")},
		{"((fn* [x y] y) a b)", tut.Sym("b-val")},
		{"((fn* (x y) y) b a)", tut.Sym("a-val")},
	}
	for _, tc := range cases {
		env := envs.MalSpecial.Nest(core.Bindings{
			"a": tut.Sym("a-val"),
			"b": tut.Sym("b-val"),
		})
		assert.Eval(t, env, tc.input, tc.want)
		assert.EnvLacks(t, env, "x", "y")
	}
}

func TestFnStarClosureCapturesDefiningEnv(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": tut.Sym("outer-a-val"),
		"d": tut.Sym("d-val"),
		"s": tut.Sym("shadow-val"),
	})
	assert.Ok(t, tut.MustEval(outer, "(def! f (fn* (x) a))"))
	assert.Eval(t, outer, "(let* (a s) (f d))", tut.Sym("outer-a-val"))
}

func TestFnStarParamShadowsOuter(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{
		"x": tut.Sym("outer-x"),
		"a": tut.Sym("a-val"),
	})
	assert.Eval(t, outer, "((fn* (x) x) a)", tut.Sym("a-val"))
	assert.EnvHas(t, outer, "x", tut.Sym("outer-x"))
}

func TestFnStarCallArityMismatch(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": tut.Sym("a-val"),
		"b": tut.Sym("b-val"),
		"c": tut.Sym("c-val"),
	})
	assert.EvalFatal(t, outer,
		"((fn* (x) x))",
		"((fn* (x) x) a b)",
		"((fn* (x y) x) a)",
		"((fn* (x y) x) a b c)",
		"((fn* [] a) a)",
		"((fn* [x] x))",
	)
}

func TestFnStarInvalidParamsAndFormArity(t *testing.T) {
	assert.EvalFatal(t, envs.MalSpecial.Nest(nil),
		"(fn* x x)",
		"(fn* a a)",
		"(fn* (()) x)",
		"(fn* [()] x)",
		"(fn* (x ()) x)",
		"(fn* [x ()] x)",
		"(fn* (x))",
		"(fn* (x) x extra)",
		"(fn* [x] x extra)",
	)
}

func TestFnStarDefIntegrationAndBodyFatal(t *testing.T) {
	env := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("a-val")})
	assert.Ok(t, tut.MustEval(env, "(def! ident (fn* (x) x))"))
	assert.Eval(t, env, "(ident a)", tut.Sym("a-val"))
	assert.Ok(t, tut.MustEval(env, "(def! ident2 (fn* [x] x))"))
	assert.Eval(t, env, "(ident2 a)", tut.Sym("a-val"))
	assert.EvalFatal(t, env, "((fn* (x) missing) a)")
}

func TestDoReturnsLastAndSequencesEffects(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": tut.Sym("a-val"),
		"b": tut.Sym("b-val"),
		"c": tut.Sym("c-val"),
	})
	assert.Eval(t, outer, "(do a)", tut.Sym("a-val"))
	assert.Eval(t, outer, "(do a b)", tut.Sym("b-val"))
	assert.Eval(t, outer, "(do a b c)", tut.Sym("c-val"))

	env := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("a-val")})
	assert.Eval(t, env, "(do (def! x a) (def! y x) y)", tut.Sym("a-val"))
	assert.EnvHas(t, env, "x", tut.Sym("a-val"))
	assert.EnvHas(t, env, "y", tut.Sym("a-val"))
}

func TestDoFatalShortCircuit(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
	assert.EvalFatal(t, outer, "(do good missing good)")
	assert.EvalFatal(t, outer, "(do missing (def! leaked good))")
	assert.EnvLacks(t, outer, "leaked")

	env := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
	assert.EvalFatal(t, env, "(do (def! x good) missing (def! y good))")
	assert.EnvHas(t, env, "x", tut.Sym("good-val"))
	assert.EnvLacks(t, env, "y")
}

func TestDoEmptyReturnsNil(t *testing.T) {
	assert.Is[typ.Nil](t, tut.MustEval(envs.MalSpecial.Nest(nil), "(do)"))
}

func TestIfTruthinessTable(t *testing.T) {
	outer := envs.MalDefault.Nest(core.Bindings{
		"t": tut.Sym("then-val"),
		"e": tut.Sym("else-val"),
	})
	assert.Eval(t, outer, "(if nil t e)", tut.Sym("else-val"))
	assert.Eval(t, outer, "(if false t e)", tut.Sym("else-val"))
	assert.Eval(t, outer, "(if true t e)", tut.Sym("then-val"))
	assert.Eval(t, outer, "(if t t e)", tut.Sym("then-val"))
	assert.Eval(t, outer, "(if e t e)", tut.Sym("then-val"))
}

func TestIfOnlyTakenBranchEvaluated(t *testing.T) {
	outer := envs.MalDefault.Nest(core.Bindings{"good": tut.Sym("good-val")})
	assert.Eval(t, outer, "(if true good missing)", tut.Sym("good-val"))
	assert.Eval(t, outer, "(if false missing good)", tut.Sym("good-val"))
	assert.Eval(t, outer, "(if nil missing good)", tut.Sym("good-val"))

	env := envs.MalDefault.Nest(core.Bindings{"good": tut.Sym("good-val")})
	assert.Eval(t, env, "(if true (def! x good) (def! y good))", tut.Sym("good-val"))
	assert.EnvHas(t, env, "x", tut.Sym("good-val"))
	assert.EnvLacks(t, env, "y")

	env2 := envs.MalDefault.Nest(core.Bindings{"good": tut.Sym("good-val")})
	assert.Eval(t, env2, "(if false (def! x good) (def! y good))", tut.Sym("good-val"))
	assert.EnvLacks(t, env2, "x")
	assert.EnvHas(t, env2, "y", tut.Sym("good-val"))
}

func TestIfCondFatalSkipsBranches(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
	assert.EvalFatal(t, outer, "(if missing good good)")

	env := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
	assert.EvalFatal(t, env, "(if missing (def! x good) (def! y good))")
	assert.EnvLacks(t, env, "x", "y")
}

func TestIfBranchFatalAndArity(t *testing.T) {
	outer := envs.MalDefault.Nest(core.Bindings{"good": tut.Sym("good-val")})
	assert.EvalFatal(t, outer,
		"(if true missing good)",
		"(if false good missing)",
		"(if nil good missing)",
		"(if true good good good)",
		"(if)",
	)
}

func TestFnStarVariadicRestCollection(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": tut.Sym("a-val"),
		"b": tut.Sym("b-val"),
		"c": tut.Sym("c-val"),
	})
	assert.Eval(t, outer, "((fn* (x & rest) x) a)", tut.Sym("a-val"))
	assert.Eval(t, outer, "((fn* (x & rest) x) a b)", tut.Sym("a-val"))
	assert.Eval(t, outer, "((fn* (x & rest) x) a b c)", tut.Sym("a-val"))
	assert.Eval(t, outer, "((fn* [x & rest] x) a b c)", tut.Sym("a-val"))

	assert.Eval(t, outer, "((fn* (x & rest) rest) a)", tut.Lst())
	assert.Eval(t, outer, "((fn* [x & rest] rest) a)", tut.Lst())
	assert.Eval(t, outer, "((fn* (a b & rest) rest) a b)", tut.Lst())

	assert.Eval(t, outer, "((fn* (x & rest) rest) a b c)",
		tut.Lst(tut.Sym("b-val"), tut.Sym("c-val")))
	assert.Eval(t, outer, "((fn* [x & rest] rest) a b c)",
		tut.Lst(tut.Sym("b-val"), tut.Sym("c-val")))

	assert.EvalFatal(t, outer,
		"((fn* (x y & rest) x) a)",
		"((fn* [x y & rest] x) a)",
		"((fn* (x & rest) x))",
	)
	assert.EnvLacks(t, outer, "x", "rest")
}

func TestFnStarVariadicMisplacedAmpersand(t *testing.T) {
	assert.EvalFatal(t, envs.MalSpecial.Nest(nil),
		"(fn* (x & ) x)",
		"(fn* [x &] x)",
		"(fn* (&) x)",
		"(fn* [&] x)",
		"(fn* (x & y z) x)",
		"(fn* [x & y z] x)",
		"(fn* (x & y & z) x)",
		"(fn* (x & ()) x)",
		"(fn* [x & ()] x)",
	)
}
