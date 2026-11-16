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
	tr := assert.New(t)
	for _, input := range []string{"(let* (x a y x) y)", "(let* [x a y x] y)"} {
		outer := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("outer-a-val")})
		tr.Eval(outer, input, tut.Sym("outer-a-val"))
		tr.EnvLacks(outer, "x", "y")
		tr.EnvHas(outer, "a", tut.Sym("outer-a-val"))
	}
}

func TestLetStarOddBindingsError(t *testing.T) {
	tr := assert.New(t)
	for _, input := range []string{"(let* (x) x)", "(let* [x] x)"} {
		_ = tr.EvalFatalContains(envs.MalSpecial.Nest(nil), input, "odd")
	}
}

func TestLetStarNonSymbolKeyQualifyPrefix(t *testing.T) {
	tr := assert.New(t)
	for _, input := range []string{"(let* (() a) a)", "(let* [() a] a)"} {
		outer := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("a-val")})
		_ = tr.EvalFatalContains(outer, input, "invalid key in let bindings")
	}
}

func TestLetStarFatalValueForwardingWithoutSet(t *testing.T) {
	tr := assert.New(t)
	for _, input := range []string{"(let* (ok missing) ok)", "(let* [ok missing] ok)"} {
		outer := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
		tr.EvalFatal(outer, input)
		tr.EnvLacks(outer, "ok")
		tr.EnvHas(outer, "good", tut.Sym("good-val"))
	}
}

func TestDefBangEvalThenSetAndReturn(t *testing.T) {
	tr := assert.New(t)
	env := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("a-val")})
	tr.Eval(env, "(def! k a)", tut.Sym("a-val"))
	tr.EnvHas(env, "k", tut.Sym("a-val"))
}

func TestDefBangFatalNonSet(t *testing.T) {
	tr := assert.New(t)
	env := envs.MalSpecial.Nest(nil)
	tr.EvalFatal(env, "(def! k missing)")
	tr.EnvLacks(env, "k")
}

func TestFnStarIdentityListVecEmptyMulti(t *testing.T) {
	tr := assert.New(t)
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
		tr.Eval(env, tc.input, tc.want)
		tr.EnvLacks(env, "x", "y")
	}
}

func TestFnStarClosureCapturesDefiningEnv(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": tut.Sym("outer-a-val"),
		"d": tut.Sym("d-val"),
		"s": tut.Sym("shadow-val"),
	})
	tr.Ok(tr.MustEval(outer, "(def! f (fn* (x) a))"))
	tr.Eval(outer, "(let* (a s) (f d))", tut.Sym("outer-a-val"))
}

func TestFnStarParamShadowsOuter(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalSpecial.Nest(core.Bindings{
		"x": tut.Sym("outer-x"),
		"a": tut.Sym("a-val"),
	})
	tr.Eval(outer, "((fn* (x) x) a)", tut.Sym("a-val"))
	tr.EnvHas(outer, "x", tut.Sym("outer-x"))
}

func TestFnStarCallArityMismatch(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": tut.Sym("a-val"),
		"b": tut.Sym("b-val"),
		"c": tut.Sym("c-val"),
	})
	tr.EvalFatal(outer,
		"((fn* (x) x))",
		"((fn* (x) x) a b)",
		"((fn* (x y) x) a)",
		"((fn* (x y) x) a b c)",
		"((fn* [] a) a)",
		"((fn* [x] x))",
	)
}

func TestFnStarInvalidParamsAndFormArity(t *testing.T) {
	tr := assert.New(t)
	tr.EvalFatal(envs.MalSpecial.Nest(nil),
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
	tr := assert.New(t)
	env := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("a-val")})
	tr.Ok(tr.MustEval(env, "(def! ident (fn* (x) x))"))
	tr.Eval(env, "(ident a)", tut.Sym("a-val"))
	tr.Ok(tr.MustEval(env, "(def! ident2 (fn* [x] x))"))
	tr.Eval(env, "(ident2 a)", tut.Sym("a-val"))
	tr.EvalFatal(env, "((fn* (x) missing) a)")
}

func TestDoReturnsLastAndSequencesEffects(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": tut.Sym("a-val"),
		"b": tut.Sym("b-val"),
		"c": tut.Sym("c-val"),
	})
	tr.Eval(outer, "(do a)", tut.Sym("a-val"))
	tr.Eval(outer, "(do a b)", tut.Sym("b-val"))
	tr.Eval(outer, "(do a b c)", tut.Sym("c-val"))

	env := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("a-val")})
	tr.Eval(env, "(do (def! x a) (def! y x) y)", tut.Sym("a-val"))
	tr.EnvHas(env, "x", tut.Sym("a-val"))
	tr.EnvHas(env, "y", tut.Sym("a-val"))
}

func TestDoFatalShortCircuit(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
	tr.EvalFatal(outer, "(do good missing good)")
	tr.EvalFatal(outer, "(do missing (def! leaked good))")
	tr.EnvLacks(outer, "leaked")

	env := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
	tr.EvalFatal(env, "(do (def! x good) missing (def! y good))")
	tr.EnvHas(env, "x", tut.Sym("good-val"))
	tr.EnvLacks(env, "y")
}

func TestDoEmptyReturnsNil(t *testing.T) {
	tr := assert.New(t)
	_ = tr.MustEvalAs[typ.Nil](envs.MalSpecial.Nest(nil), "(do)")
}

func TestIfTruthinessTable(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalDefault.Nest(core.Bindings{
		"t": tut.Sym("then-val"),
		"e": tut.Sym("else-val"),
	})
	tr.Eval(outer, "(if nil t e)", tut.Sym("else-val"))
	tr.Eval(outer, "(if false t e)", tut.Sym("else-val"))
	tr.Eval(outer, "(if true t e)", tut.Sym("then-val"))
	tr.Eval(outer, "(if t t e)", tut.Sym("then-val"))
	tr.Eval(outer, "(if e t e)", tut.Sym("then-val"))
}

func TestIfOnlyTakenBranchEvaluated(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalDefault.Nest(core.Bindings{"good": tut.Sym("good-val")})
	tr.Eval(outer, "(if true good missing)", tut.Sym("good-val"))
	tr.Eval(outer, "(if false missing good)", tut.Sym("good-val"))
	tr.Eval(outer, "(if nil missing good)", tut.Sym("good-val"))

	env := envs.MalDefault.Nest(core.Bindings{"good": tut.Sym("good-val")})
	tr.Eval(env, "(if true (def! x good) (def! y good))", tut.Sym("good-val"))
	tr.EnvHas(env, "x", tut.Sym("good-val"))
	tr.EnvLacks(env, "y")

	env2 := envs.MalDefault.Nest(core.Bindings{"good": tut.Sym("good-val")})
	tr.Eval(env2, "(if false (def! x good) (def! y good))", tut.Sym("good-val"))
	tr.EnvLacks(env2, "x")
	tr.EnvHas(env2, "y", tut.Sym("good-val"))
}

func TestIfCondFatalSkipsBranches(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
	tr.EvalFatal(outer, "(if missing good good)")

	env := envs.MalSpecial.Nest(core.Bindings{"good": tut.Sym("good-val")})
	tr.EvalFatal(env, "(if missing (def! x good) (def! y good))")
	tr.EnvLacks(env, "x", "y")
}

func TestIfBranchFatalAndArity(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalDefault.Nest(core.Bindings{"good": tut.Sym("good-val")})
	tr.EvalFatal(outer,
		"(if true missing good)",
		"(if false good missing)",
		"(if nil good missing)",
		"(if true good good good)",
		"(if)",
	)
}

func TestFnStarVariadicRestCollection(t *testing.T) {
	tr := assert.New(t)
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": tut.Sym("a-val"),
		"b": tut.Sym("b-val"),
		"c": tut.Sym("c-val"),
	})
	tr.Eval(outer, "((fn* (x & rest) x) a)", tut.Sym("a-val"))
	tr.Eval(outer, "((fn* (x & rest) x) a b)", tut.Sym("a-val"))
	tr.Eval(outer, "((fn* (x & rest) x) a b c)", tut.Sym("a-val"))
	tr.Eval(outer, "((fn* [x & rest] x) a b c)", tut.Sym("a-val"))

	tr.Eval(outer, "((fn* (x & rest) rest) a)", tut.Lst())
	tr.Eval(outer, "((fn* [x & rest] rest) a)", tut.Lst())
	tr.Eval(outer, "((fn* (a b & rest) rest) a b)", tut.Lst())

	tr.Eval(outer, "((fn* (x & rest) rest) a b c)",
		tut.Lst(tut.Sym("b-val"), tut.Sym("c-val")))
	tr.Eval(outer, "((fn* [x & rest] rest) a b c)",
		tut.Lst(tut.Sym("b-val"), tut.Sym("c-val")))

	tr.EvalFatal(outer,
		"((fn* (x y & rest) x) a)",
		"((fn* [x y & rest] x) a)",
		"((fn* (x & rest) x))",
	)
	tr.EnvLacks(outer, "x", "rest")
}

func TestFnStarVariadicMisplacedAmpersand(t *testing.T) {
	tr := assert.New(t)
	tr.EvalFatal(envs.MalSpecial.Nest(nil),
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

func TestQuoteSymbolUnevaluated(t *testing.T) {
	tr := assert.New(t)
	env := envs.MalSpecial.Nest(core.Bindings{"a": tut.Sym("a-val")})
	tr.Eval(env, "(quote a)", tut.Sym("a"))
}

func TestQuoteListUnevaluated(t *testing.T) {
	tr := assert.New(t)
	tr.Eval(envs.MalSpecial.Nest(nil), "(quote (a b))", tut.Lst(tut.Sym("a"), tut.Sym("b")))
}

func TestQuoteEval(t *testing.T) {
	tr := assert.New(t)
	envr := envs.MalSpecial.Nest(nil)

	tr.Equal(tr.MustEvalAs[typ.Integer](envr, "'42"), typ.Integer(42))
	_ = tr.MustEvalAs[typ.Nil](envr, "'nil")
	tr.Equal(tr.MustEvalAs[typ.String](envr, `'"hi"`), typ.String("hi"))

	// Quote type itself: String, Type and Eval return raw value.
	q := tut.Quote(typ.Integer(1))
	tr.Repr[typ.List](q, "(quote 1)")
	tr.EvalValue(envr, q, typ.Integer(1))
}

func TestSymbolViaQuoteMustEval(t *testing.T) {
	tr := assert.New(t)
	envr := envs.MalSpecial.Nest(nil)

	// Quoting prevents symbol lookup, so this works with an env that does not contain my-sym and
	// proves Quote returns the raw unevaluated form.
	got := tr.MustEvalAs[typ.Symbol](envr, "'my-sym")
	tr.Equal(string(got), "my-sym")
	tr.Repr[typ.Symbol](got, "my-sym")

	// Direct Symbol trivial behaviour (no env lookup).
	direct := typ.Symbol("foo")
	tr.Repr[typ.Symbol](direct, "foo")
}
