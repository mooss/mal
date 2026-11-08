package envs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func fatalMsg(v core.Value) string {
	return fmt.Sprintf("%v", v)
}

func TestLetStarSequentialScopingAndOuterIsolation(t *testing.T) {
	outer := core.NewEnv(core.Bindings{
		"a": typ.Symbol("outer-a-val"),
	})
	bindings := typ.List{
		typ.Symbol("x"), typ.Symbol("a"),
		typ.Symbol("y"), typ.Symbol("x"),
	}
	form := typ.Symbol("y")

	res := letStar(outer, bindings, form)
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

func TestLetStarOddBindingsError(t *testing.T) {
	outer := core.NewEnv()
	res := letStar(outer, typ.List{typ.Symbol("x")}, typ.Symbol("x"))
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected fatal for odd bindings, got %v", res)
	}
	if msg := fatalMsg(res); !strings.Contains(msg, "odd") {
		t.Fatalf("expected odd-numbered error, got %q", msg)
	}
}

func TestLetStarNonSymbolKeyQualifyPrefix(t *testing.T) {
	outer := core.NewEnv(core.Bindings{
		"a": typ.Symbol("a-val"),
	})
	// empty list is a value but not a symbol.
	bindings := typ.List{
		typ.List{}, typ.Symbol("a"),
	}
	res := letStar(outer, bindings, typ.Symbol("a"))
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected fatal for non-symbol key, got %v", res)
	}
	if msg := fatalMsg(res); !strings.Contains(msg, "invalid key in let bindings") {
		t.Fatalf("expected qualify prefix %q, got %q", "invalid key in let bindings", msg)
	}
}

func TestLetStarFatalValueForwardingWithoutSet(t *testing.T) {
	outer := core.NewEnv(core.Bindings{
		"good": typ.Symbol("good-val"),
	})
	bindings := typ.List{
		typ.Symbol("ok"), typ.Symbol("missing"),
	}
	res := letStar(outer, bindings, typ.Symbol("ok"))
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

func TestDefBangEvalThenSetAndReturn(t *testing.T) {
	env := core.NewEnv(core.Bindings{
		"a": typ.Symbol("a-val"),
	})
	res := defBang(env, typ.Symbol("k"), typ.Symbol("a"))
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
	env := core.NewEnv()
	res := defBang(env, typ.Symbol("k"), typ.Symbol("missing"))
	if !ut.Is[typ.Fatal](res) {
		t.Fatalf("expected fatal, got %v", res)
	}
	if v := env.Get("k"); v != nil {
		t.Fatalf("fatal def! should not set, got %v", v)
	}
}
