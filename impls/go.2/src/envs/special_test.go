//revive:disable:cognitive-complexity,cyclomatic
//nolint:dupword
package envs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func fatalMsg(v core.Value) string {
	return fmt.Sprintf("%v", v)
}

func TestLetStarSequentialScopingAndOuterIsolation(t *testing.T) {
	for _, input := range []string{"(let* (x a y x) y)", "(let* [x a y x] y)"} {
		outer := envs.MalSpecial.Nest(core.Bindings{"a": typ.Symbol("outer-a-val")})
		res := tut.MustEval(outer, input)
		if ut.Is[typ.Fatal](res) {
			t.Fatalf("expected success, got fatal: %v", res)
		}
		if res != typ.Symbol("outer-a-val") {
			t.Fatalf("expected sequential scoping result %q, got %v", "outer-a-val", res)
		}
		if v := outer.Get("x"); v != nil {
			t.Fatalf("outer env leaked x: %v", v)
		}
		if v := outer.Get("y"); v != nil {
			t.Fatalf("outer env leaked y: %v", v)
		}
		if v := outer.Get("a"); v != typ.Symbol("outer-a-val") {
			t.Fatalf("outer binding changed: %v", v)
		}
	}
}

func TestLetStarOddBindingsError(t *testing.T) {
	for _, input := range []string{"(let* (x) x)", "(let* [x] x)"} {
		res := tut.MustEval(envs.MalSpecial.Nest(nil), input)
		if !ut.Is[typ.Fatal](res) {
			t.Fatalf("expected fatal for odd bindings, got %v", res)
		}
		if msg := fatalMsg(res); !strings.Contains(msg, "odd") {
			t.Fatalf("expected odd-numbered error, got %q", msg)
		}
	}
}

//nolint:dupword
func TestLetStarNonSymbolKeyQualifyPrefix(t *testing.T) {
	for _, input := range []string{"(let* (() a) a)", "(let* [() a] a)"} {
		outer := envs.MalSpecial.Nest(core.Bindings{"a": typ.Symbol("a-val")})
		// empty list is a value but not a symbol.
		res := tut.MustEval(outer, input)
		if !ut.Is[typ.Fatal](res) {
			t.Fatalf("expected fatal for non-symbol key, got %v", res)
		}
		if msg := fatalMsg(res); !strings.Contains(msg, "invalid key in let bindings") {
			t.Fatalf("expected qualify prefix %q, got %q", "invalid key in let bindings", msg)
		}
	}
}

func TestLetStarFatalValueForwardingWithoutSet(t *testing.T) {
	for _, input := range []string{"(let* (ok missing) ok)", "(let* [ok missing] ok)"} {
		outer := envs.MalSpecial.Nest(core.Bindings{"good": typ.Symbol("good-val")})
		res := tut.MustEval(outer, input)
		if !ut.Is[typ.Fatal](res) {
			t.Fatalf("expected fatal forwarding, got %v", res)
		}
		if v := outer.Get("ok"); v != nil {
			t.Fatalf("fatal binding should not leak to outer, got %v", v)
		}
		if v := outer.Get("good"); v != typ.Symbol("good-val") {
			t.Fatalf("outer binding changed: %v", v)
		}
	}
}

func TestDefBangEvalThenSetAndReturn(t *testing.T) {
	env := envs.MalSpecial.Nest(core.Bindings{"a": typ.Symbol("a-val")})
	res := tut.MustEval(env, "(def! k a)")
	if ut.Is[typ.Fatal](res) {
		t.Fatalf("expected success, got fatal: %v", res)
	}
	if res != typ.Symbol("a-val") {
		t.Fatalf("expected evaled return %q, got %v", "a-val", res)
	}
	if v := env.Get("k"); v != typ.Symbol("a-val") {
		t.Fatalf("expected env k=%q, got %v", "a-val", v)
	}
}

func TestDefBangFatalNonSet(t *testing.T) {
	env := envs.MalSpecial.Nest(nil)
	res := tut.MustEval(env, "(def! k missing)")
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected fatal, got %v", res)
	}
	if v := env.Get("k"); v != nil {
		t.Fatalf("fatal def! should not set, got %v", v)
	}
}

func TestFnStarIdentityListVecEmptyMulti(t *testing.T) {
	cases := []struct {
		input string
		want  typ.Symbol
	}{
		{"((fn* (x) x) a)", typ.Symbol("a-val")},
		{"((fn* [x] x) a)", typ.Symbol("a-val")},
		{"((fn* () a))", typ.Symbol("a-val")},
		{"((fn* [] a))", typ.Symbol("a-val")},
		{"((fn* (x y) x) a b)", typ.Symbol("a-val")},
		{"((fn* [x y] y) a b)", typ.Symbol("b-val")},
		{"((fn* (x y) y) b a)", typ.Symbol("a-val")},
	}
	for _, tc := range cases {
		env := envs.MalSpecial.Nest(core.Bindings{
			"a": typ.Symbol("a-val"),
			"b": typ.Symbol("b-val"),
		})
		res := tut.MustEval(env, tc.input)
		if ut.Is[typ.Fatal](res) {
			t.Fatalf("fn* %q expected success, got fatal: %v", tc.input, res)
		}
		if res != tc.want {
			t.Fatalf("fn* %q expected %v, got %v", tc.input, tc.want, res)
		}
		if v := env.Get("x"); v != nil {
			t.Fatalf("fn* %q leaked x: %v", tc.input, v)
		}
		if v := env.Get("y"); v != nil {
			t.Fatalf("fn* %q leaked y: %v", tc.input, v)
		}
	}
}

func TestFnStarClosureCapturesDefiningEnv(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": typ.Symbol("outer-a-val"),
		"d": typ.Symbol("d-val"),
		"s": typ.Symbol("shadow-val"),
	})
	if res := tut.MustEval(outer, "(def! f (fn* (x) a))"); ut.Is[typ.Fatal](res) {
		t.Fatalf("def f failed: %v", res)
	}
	res := tut.MustEval(outer, "(let* (a s) (f d))")
	if ut.Is[typ.Fatal](res) {
		t.Fatalf("expected lexical capture, got fatal: %v", res)
	}
	if res != typ.Symbol("outer-a-val") {
		t.Fatalf("expected closure to see defining env %v, got %v", typ.Symbol("outer-a-val"), res)
	}
}

func TestFnStarParamShadowsOuter(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{
		"x": typ.Symbol("outer-x"),
		"a": typ.Symbol("a-val"),
	})
	res := tut.MustEval(outer, "((fn* (x) x) a)")
	if ut.Is[typ.Fatal](res) {
		t.Fatalf("expected success, got fatal: %v", res)
	}
	if res != typ.Symbol("a-val") {
		t.Fatalf("expected param shadowing %v, got %v", typ.Symbol("a-val"), res)
	}
	if v := outer.Get("x"); v != typ.Symbol("outer-x") {
		t.Fatalf("outer x changed: %v", v)
	}
}

func TestFnStarCallArityMismatch(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": typ.Symbol("a-val"),
		"b": typ.Symbol("b-val"),
		"c": typ.Symbol("c-val"),
	})
	for _, input := range []string{
		"((fn* (x) x))",
		"((fn* (x) x) a b)",
		"((fn* (x y) x) a)",
		"((fn* (x y) x) a b c)",
		"((fn* [] a) a)",
		"((fn* [x] x))",
	} {
		res := tut.MustEval(outer, input)
		if !ut.Is[typ.Fatal](res) {
			t.Fatalf("expected arity fatal for %q, got %v", input, res)
		}
	}
}

func TestFnStarInvalidParamsAndFormArity(t *testing.T) {
	for _, input := range []string{
		"(fn* x x)",
		"(fn* a a)",
	} {
		res := tut.MustEval(envs.MalSpecial.Nest(nil), input)
		if !ut.Is[typ.Fatal](res) {
			t.Fatalf("expected fatal for non-list/vec params %q, got %v", input, res)
		}
	}
	for _, input := range []string{
		"(fn* (()) x)",
		"(fn* [()] x)",
		"(fn* (x ()) x)",
		"(fn* [x ()] x)",
	} {
		res := tut.MustEval(envs.MalSpecial.Nest(nil), input)
		if !ut.Is[typ.Fatal](res) {
			t.Fatalf("expected fatal for non-symbol param %q, got %v", input, res)
		}
	}
	for _, input := range []string{
		"(fn* (x))",
		"(fn* (x) x extra)",
		"(fn* [x] x extra)",
	} {
		res := tut.MustEval(envs.MalSpecial.Nest(nil), input)
		if !ut.Is[typ.Fatal](res) {
			t.Fatalf("expected form-arity fatal for %q, got %v", input, res)
		}
	}
}

func TestFnStarDefIntegrationAndBodyFatal(t *testing.T) {
	env := envs.MalSpecial.Nest(core.Bindings{"a": typ.Symbol("a-val")})
	if res := tut.MustEval(env, "(def! ident (fn* (x) x))"); ut.Is[typ.Fatal](res) {
		t.Fatalf("def ident failed: %v", res)
	}
	res := tut.MustEval(env, "(ident a)")
	if ut.Is[typ.Fatal](res) {
		t.Fatalf("ident call failed: %v", res)
	}
	if res != typ.Symbol("a-val") {
		t.Fatalf("expected ident %v, got %v", typ.Symbol("a-val"), res)
	}
	if res := tut.MustEval(env, "(def! ident2 (fn* [x] x))"); ut.Is[typ.Fatal](res) {
		t.Fatalf("def ident2 failed: %v", res)
	}
	res = tut.MustEval(env, "(ident2 a)")
	if res != typ.Symbol("a-val") {
		t.Fatalf("expected ident2 %v, got %v", typ.Symbol("a-val"), res)
	}
	res = tut.MustEval(env, "((fn* (x) missing) a)")
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected body fatal, got %v", res)
	}
}

func TestDoReturnsLastAndSequencesEffects(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{
		"a": typ.Symbol("a-val"),
		"b": typ.Symbol("b-val"),
		"c": typ.Symbol("c-val"),
	})
	cases := []struct {
		input string
		want  typ.Symbol
	}{
		{"(do a)", typ.Symbol("a-val")},
		{"(do a b)", typ.Symbol("b-val")},
		{"(do a b c)", typ.Symbol("c-val")},
	}
	for _, tc := range cases {
		res := tut.MustEval(outer, tc.input)
		if ut.Is[typ.Fatal](res) {
			t.Fatalf("do %q expected success, got fatal: %v", tc.input, res)
		}
		if res != tc.want {
			t.Fatalf("do %q expected %v, got %v", tc.input, tc.want, res)
		}
	}
	env := envs.MalSpecial.Nest(core.Bindings{"a": typ.Symbol("a-val")})
	res := tut.MustEval(env, "(do (def! x a) (def! y x) y)")
	if ut.Is[typ.Fatal](res) {
		t.Fatalf("do sequence failed: %v", res)
	}
	if res != typ.Symbol("a-val") {
		t.Fatalf("expected do sequence %v, got %v", typ.Symbol("a-val"), res)
	}
	if v := env.Get("x"); v != typ.Symbol("a-val") {
		t.Fatalf("expected x=%v, got %v", typ.Symbol("a-val"), v)
	}
	if v := env.Get("y"); v != typ.Symbol("a-val") {
		t.Fatalf("expected y=%v, got %v", typ.Symbol("a-val"), v)
	}
}

func TestDoFatalShortCircuit(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{"good": typ.Symbol("good-val")})
	res := tut.MustEval(outer, "(do good missing good)")
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected fatal, got %v", res)
	}
	res = tut.MustEval(outer, "(do missing (def! leaked good))")
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected fatal, got %v", res)
	}
	if v := outer.Get("leaked"); v != nil {
		t.Fatalf("do should short-circuit, leaked set: %v", v)
	}
	env := envs.MalSpecial.Nest(core.Bindings{"good": typ.Symbol("good-val")})
	res = tut.MustEval(env, "(do (def! x good) missing (def! y good))")
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected fatal, got %v", res)
	}
	if v := env.Get("x"); v != typ.Symbol("good-val") {
		t.Fatalf("expected x set before fatal, got %v", v)
	}
	if v := env.Get("y"); v != nil {
		t.Fatalf("do should skip after fatal, y set: %v", v)
	}
}

func TestDoEmptyReturnsNil(t *testing.T) {
	res := tut.MustEval(envs.MalSpecial.Nest(nil), "(do)")
	if res != nil {
		if _, ok := res.(typ.Nil); !ok {
			t.Fatalf("expected nil for empty do, got %v (%T)", res, res)
		}
	}
}

func TestIfTruthinessTable(t *testing.T) {
	outer := envs.MalDefault.Nest(core.Bindings{
		"t": typ.Symbol("then-val"),
		"e": typ.Symbol("else-val"),
	})
	cases := []struct {
		input string
		want  typ.Symbol
	}{
		{"(if nil t e)", typ.Symbol("else-val")},
		{"(if false t e)", typ.Symbol("else-val")},
		{"(if true t e)", typ.Symbol("then-val")},
		{"(if t t e)", typ.Symbol("then-val")},
		{"(if e t e)", typ.Symbol("then-val")},
	}
	for _, tc := range cases {
		res := tut.MustEval(outer, tc.input)
		if ut.Is[typ.Fatal](res) {
			t.Fatalf("if %q expected success, got fatal: %v", tc.input, res)
		}
		if res != tc.want {
			t.Fatalf("if %q expected %v, got %v", tc.input, tc.want, res)
		}
	}
}

func TestIfOnlyTakenBranchEvaluated(t *testing.T) {
	outer := envs.MalDefault.Nest(core.Bindings{"good": typ.Symbol("good-val")})
	for _, input := range []string{
		"(if true good missing)",
		"(if false missing good)",
		"(if nil missing good)",
	} {
		res := tut.MustEval(outer, input)
		if ut.Is[typ.Fatal](res) {
			t.Fatalf("if %q expected success, got fatal: %v", input, res)
		}
		if res != typ.Symbol("good-val") {
			t.Fatalf("if %q expected %v, got %v", input, typ.Symbol("good-val"), res)
		}
	}
	env := envs.MalDefault.Nest(core.Bindings{"good": typ.Symbol("good-val")})
	res := tut.MustEval(env, "(if true (def! x good) (def! y good))")
	if ut.Is[typ.Fatal](res) {
		t.Fatalf("if-true def failed: %v", res)
	}
	if res != typ.Symbol("good-val") {
		t.Fatalf("expected %v, got %v", typ.Symbol("good-val"), res)
	}
	if v := env.Get("x"); v != typ.Symbol("good-val") {
		t.Fatalf("then branch should set x, got %v", v)
	}
	if v := env.Get("y"); v != nil {
		t.Fatalf("else branch should not run, y set: %v", v)
	}
	env2 := envs.MalDefault.Nest(core.Bindings{"good": typ.Symbol("good-val")})
	res = tut.MustEval(env2, "(if false (def! x good) (def! y good))")
	if ut.Is[typ.Fatal](res) {
		t.Fatalf("if-false def failed: %v", res)
	}
	if v := env2.Get("x"); v != nil {
		t.Fatalf("then branch should not run, x set: %v", v)
	}
	if v := env2.Get("y"); v != typ.Symbol("good-val") {
		t.Fatalf("else branch should set y, got %v", v)
	}
}

func TestIfCondFatalSkipsBranches(t *testing.T) {
	outer := envs.MalSpecial.Nest(core.Bindings{"good": typ.Symbol("good-val")})
	res := tut.MustEval(outer, "(if missing good good)")
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected cond fatal, got %v", res)
	}
	env := envs.MalSpecial.Nest(core.Bindings{"good": typ.Symbol("good-val")})
	res = tut.MustEval(env, "(if missing (def! x good) (def! y good))")
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected cond fatal, got %v", res)
	}
	if v := env.Get("x"); v != nil {
		t.Fatalf("cond fatal should skip then, x set: %v", v)
	}
	if v := env.Get("y"); v != nil {
		t.Fatalf("cond fatal should skip else, y set: %v", v)
	}
}

func TestIfBranchFatalAndArity(t *testing.T) {
	outer := envs.MalDefault.Nest(core.Bindings{"good": typ.Symbol("good-val")})
	for _, input := range []string{
		"(if true missing good)",
		"(if false good missing)",
		"(if nil good missing)",
	} {
		res := tut.MustEval(outer, input)
		if !ut.Is[typ.Fatal](res) {
			t.Fatalf("expected branch fatal for %q, got %v", input, res)
		}
	}
	for _, input := range []string{
		"(if true good good good)",
		"(if)",
	} {
		res := tut.MustEval(outer, input)
		if !ut.Is[typ.Fatal](res) {
			t.Fatalf("expected arity fatal for %q, got %v", input, res)
		}
	}
}
