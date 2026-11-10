package core_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
)

type _tv string

func (_tv) Type() core.TypeName           { return "test" }
func (v _tv) String() string              { return string(v) }
func (v _tv) Eval(core.Env) core.Value    { return v }
func (v _tv) Equal(other core.Value) bool { return v == other }

func TestNewEnvMergesAndCopiesBindings(t *testing.T) {
	first := core.Bindings{
		"first":  _tv("first"),
		"shared": _tv("first"),
	}
	second := core.Bindings{
		"second": _tv("second"),
		"shared": _tv("second"),
	}

	env := core.NewEnv(first, second)
	for key, want := range map[string]core.Value{
		"first":  _tv("first"),
		"second": _tv("second"),
		"shared": _tv("second"),
	} {
		if got := env.Get(key); got != want {
			t.Errorf("Get(%q) = %v, want %v", key, got, want)
		}
	}

	first["first"] = _tv("changed")
	if got := env.Get("first"); got != _tv("first") {
		msg := "Get(%q) after changing source bindings = %v, want %v"
		t.Errorf(msg, "first", got, _tv("first"))
	}
}

func TestEnvGetSearchesOuterScopes(t *testing.T) {
	grandparent := core.NewEnv(core.Bindings{"grandparent": _tv("grandparent")})
	parent := grandparent.Nest(core.Bindings{"parent": _tv("parent"), "shadowed": _tv("parent")})
	child := parent.Nest(core.Bindings{"local": _tv("local"), "shadowed": _tv("local")})

	for key, want := range map[string]core.Value{
		"local":       _tv("local"),
		"parent":      _tv("parent"),
		"grandparent": _tv("grandparent"),
		"shadowed":    _tv("local"),
	} {
		if got := child.Get(key); got != want {
			t.Errorf("Get(%q) = %v, want %v", key, got, want)
		}
	}
	if got := child.Get("missing"); got != nil {
		t.Errorf("Get(%q) = %v, want nil", "missing", got)
	}
}

func TestEnvSetUpdatesOnlyCurrentScope(t *testing.T) {
	parent := core.NewEnv(core.Bindings{
		"shared": _tv("parent"),
		"parent": _tv("parent"),
	})

	env := parent.Nest(core.Bindings{})

	env.Set("shared", _tv("child"))
	env.Set("local", _tv("local"))

	if got := env.Get("shared"); got != _tv("child") {
		t.Errorf("Get(%q) = %v, want %v", "shared", got, _tv("child"))
	}
	if got := parent.Get("shared"); got != _tv("parent") {
		t.Errorf("parent Get(%q) = %v, want %v", "shared", got, _tv("parent"))
	}
	if got := env.Get("local"); got != _tv("local") {
		t.Errorf("inner[%q] = %v, want %v", "local", got, _tv("local"))
	}
}

func TestBindingsAsEnvSharesBindings(t *testing.T) {
	bindings := core.Bindings{"key": _tv("before")}
	env := bindings.AsEnv()

	bindings["key"] = _tv("after")
	if got := env.Get("key"); got != _tv("after") {
		t.Errorf("Get(%q) after changing bindings = %v, want %v", "key", got, _tv("after"))
	}

	env.Set("key", _tv("set"))
	if got := bindings["key"]; got != _tv("set") {
		t.Errorf("bindings[%q] after Set = %v, want %v", "key", got, _tv("set"))
	}
}

func TestEnvNestLookupShadowsAndFallsBack(t *testing.T) {
	outer := core.NewEnv(core.Bindings{
		"outer":  _tv("outer"),
		"shared": _tv("outer"),
	})
	nested := outer.Nest(core.Bindings{
		"inner":  _tv("inner"),
		"shared": _tv("inner"),
	})

	for key, want := range map[string]core.Value{
		"inner":  _tv("inner"),
		"outer":  _tv("outer"),
		"shared": _tv("inner"),
	} {
		if got := nested.Get(key); got != want {
			t.Errorf("Get(%q) = %v, want %v", key, got, want)
		}
	}
	if got := nested.Get("missing"); got != nil {
		t.Errorf("Get(%q) = %v, want nil", "missing", got)
	}
	if got := outer.Get("inner"); got != nil {
		t.Errorf("outer Get(%q) = %v, want nil", "inner", got)
	}
	if got := outer.Get("shared"); got != _tv("outer") {
		t.Errorf("outer Get(%q) = %v, want %v", "shared", got, _tv("outer"))
	}

	fallback := outer.Nest(nil)
	if got := fallback.Get("outer"); got != _tv("outer") {
		t.Errorf("Get(%q) with nil inner = %v, want %v", "outer", got, _tv("outer"))
	}
	if got := fallback.Get("missing"); got != nil {
		t.Errorf("Get(%q) with nil inner = %v, want nil", "missing", got)
	}
}

func TestCatEnvRightmostAndInnermostPrecedenceFlat(t *testing.T) {
	left := core.NewEnv(core.Bindings{
		"left":   _tv("left"),
		"shared": _tv("left"),
	})
	rightOuter := core.NewEnv(core.Bindings{
		"right-outer": _tv("right-outer"),
		"shared":      _tv("right-outer"),
		"nested":      _tv("outer"),
	})
	right := rightOuter.Nest(core.Bindings{
		"right-inner": _tv("right-inner"),
		"shared":      _tv("right-inner"),
		"nested":      _tv("inner"),
	})

	merged := core.CatEnv(left, right)
	for key, want := range map[string]core.Value{
		"left":        _tv("left"),
		"right-outer": _tv("right-outer"),
		"right-inner": _tv("right-inner"),
		"shared":      _tv("right-inner"),
		"nested":      _tv("inner"),
	} {
		if got := merged.Get(key); got != want {
			t.Errorf("Get(%q) = %v, want %v", key, got, want)
		}
	}

	leftOnly := core.CatEnv(left)
	if got := leftOnly.Get("shared"); got != _tv("left") {
		t.Errorf("Get(%q) = %v, want %v", "shared", got, _tv("left"))
	}

	merged.Set("left", _tv("changed"))
	if got := left.Get("left"); got != _tv("left") {
		t.Errorf("left Get(%q) after changing merged = %v, want %v", "left", got, _tv("left"))
	}
}

func TestEnvGetMissingReturnsNil(t *testing.T) {
	if got := core.NewEnv().Get("missing"); got != nil {
		t.Errorf("Get(%q) on empty env = %v, want nil", "missing", got)
	}
	env := core.NewEnv(core.Bindings{"present": _tv("present")}).Nest(core.Bindings{
		"local": _tv("local"),
	})
	if got := env.Get("missing"); got != nil {
		t.Errorf("Get(%q) on nested env = %v, want nil", "missing", got)
	}
}

func TestBindingsAsEnvNilAllowsSet(t *testing.T) {
	var bindings core.Bindings
	env := bindings.AsEnv()
	env.Set("key", _tv("value"))
	if got := env.Get("key"); got != _tv("value") {
		t.Errorf("Get(%q) = %v, want %v", "key", got, _tv("value"))
	}
}
