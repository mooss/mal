package typ

import (
	"strings"
	"testing"

	"github.com/mooss/mal/src/core"
)

func mkTestSpecial() Special {
	return Special(func(core.Env, ...core.Value) core.Value { return Fatal{} })
}

func TestSpIIRejectsWrongArgCount(t *testing.T) {
	var env core.Env
	called := false
	form := SpII(func(_ core.Env, a Special, _ Special) Special {
		called = true
		return a
	})

	s := mkTestSpecial()
	cases := [][]core.Value{
		{},
		{s},
		{s, s, s},
		{s, s, s, s},
	}
	for i, args := range cases {
		called = false
		res := form(env, args...)
		f, ok := res.(Fatal)
		if !ok {
			t.Fatalf("case %d: expected Fatal, got %T (%v)", i, res, res)
		}
		if !f.NotNil() {
			t.Fatalf("case %d: expected NotNil fatal for arg count %d", i, len(args))
		}
		if called {
			t.Fatalf("case %d: fun should not be called with %d args", i, len(args))
		}
		if !strings.Contains(strings.ToLower(f.String()), "expected 2 args") &&
			!strings.Contains(f.String(), "expected 2") {
			t.Fatalf("case %d: expected arg count message, got %q", i, f.String())
		}
	}
}

func TestSpIIForwardsTypeMismatch(t *testing.T) {
	var env core.Env
	called := false
	form := SpII(func(core.Env, Fatal, Fatal) Fatal {
		called = true
		return Fatal{}
	})

	s := mkTestSpecial()
	// first arg mismatch
	res := form(env, s, Fatalf("ok"))
	f, ok := res.(Fatal)
	if !ok {
		t.Fatalf("expected Fatal for type mismatch, got %T", res)
	}
	if !f.NotNil() {
		t.Fatal("expected NotNil fatal for type mismatch")
	}
	if called {
		t.Fatal("fun should not be called on type mismatch")
	}
	if !strings.Contains(strings.ToLower(f.String()), "expected") {
		t.Fatalf("expected mismatch message to contain 'expected', got %q", f.String())
	}

	// second arg mismatch
	called = false
	res = form(env, Fatalf("ok"), s)
	f, ok = res.(Fatal)
	if !ok {
		t.Fatalf("expected Fatal for second arg mismatch, got %T", res)
	}
	if !f.NotNil() {
		t.Fatal("expected NotNil fatal for second arg mismatch")
	}
	if called {
		t.Fatal("fun should not be called on second arg mismatch")
	}
	if !strings.Contains(strings.ToLower(f.String()), "expected") {
		t.Fatalf("expected mismatch message to contain 'expected', got %q", f.String())
	}
}

func TestSpIIForwardsIncomingFatal(t *testing.T) {
	var env core.Env
	called := false
	form := SpII(func(_ core.Env, a Special, _ Special) Special {
		called = true
		return a
	})

	s := mkTestSpecial()
	incoming := Fatalf("boom %d", 42)

	res := form(env, incoming, s)
	f, ok := res.(Fatal)
	if !ok {
		t.Fatalf("expected Fatal forwarding, got %T", res)
	}
	if !f.NotNil() {
		t.Fatal("forwarded fatal should be NotNil")
	}
	if called {
		t.Fatal("fun should not be called when incoming fatal")
	}
	if f.String() != incoming.String() {
		t.Fatalf("expected forwarded %q, got %q", incoming.String(), f.String())
	}

	called = false
	res = form(env, s, incoming)
	f, ok = res.(Fatal)
	if !ok {
		t.Fatalf("expected Fatal forwarding second arg, got %T", res)
	}
	if f.String() != incoming.String() {
		t.Fatalf("expected forwarded %q, got %q", incoming.String(), f.String())
	}
	if called {
		t.Fatal("fun should not be called when second arg is fatal")
	}
}

func TestSpIIPassesCorrectTypedArgs(t *testing.T) {
	var env core.Env
	var gotA, gotB Fatal
	form := SpII(func(_ core.Env, a Fatal, b Fatal) Fatal {
		gotA = a
		gotB = b
		return Fatalf("ok-result")
	})

	a := Fatalf("first")
	b := Fatalf("second")
	res := form(env, a, b)
	f, ok := res.(Fatal)
	if !ok {
		t.Fatalf("expected Fatal result, got %T", res)
	}
	if !strings.Contains(f.String(), "ok-result") {
		t.Fatalf("expected result passthrough, got %q", f.String())
	}
	if gotA.String() != a.String() {
		t.Fatalf("first arg not passed correctly: want %q got %q", a.String(), gotA.String())
	}
	if gotB.String() != b.String() {
		t.Fatalf("second arg not passed correctly: want %q got %q", b.String(), gotB.String())
	}
}

func TestQualifyPrefixesMessages(t *testing.T) {
	base := Fatalf("base")
	q := base.Qualify("prefix: ")
	if !q.NotNil() {
		t.Fatal("qualified fatal should be NotNil")
	}
	s := q.String()
	if !strings.Contains(s, "prefix: ") {
		t.Fatalf("qualified message should contain prefix, got %q", s)
	}
	if !strings.Contains(s, "base") {
		t.Fatalf("qualified message should contain base, got %q", s)
	}
}

func TestNotNilDistinguishesEmptyVersusFilledFatal(t *testing.T) {
	var empty Fatal
	if empty.NotNil() {
		t.Fatal("empty Fatal{} should not be NotNil")
	}
	filled := Fatalf("filled")
	if !filled.NotNil() {
		t.Fatal("Fatalf fatal should be NotNil")
	}
}
