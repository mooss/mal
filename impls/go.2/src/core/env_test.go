package core

import "testing"

type envTestValue string

func (envTestValue) Type() TypeName   { return "test" }
func (v envTestValue) String() string { return string(v) }
func (v envTestValue) Eval(Env) Value { return v }

func TestNewEnvMergesAndCopiesBindings(t *testing.T) {
	first := Bindings{
		"first":  envTestValue("first"),
		"shared": envTestValue("first"),
	}
	second := Bindings{
		"second": envTestValue("second"),
		"shared": envTestValue("second"),
	}

	env := NewEnv(first, second)
	for key, want := range map[string]Value{
		"first":  envTestValue("first"),
		"second": envTestValue("second"),
		"shared": envTestValue("second"),
	} {
		if got := env.Get(key); got != want {
			t.Errorf("Get(%q) = %v, want %v", key, got, want)
		}
	}

	first["first"] = envTestValue("changed")
	if got := env.Get("first"); got != envTestValue("first") {
		msg := "Get(%q) after changing source bindings = %v, want %v"
		t.Errorf(msg, "first", got, envTestValue("first"))
	}
}

func TestEnvGetSearchesOuterScopes(t *testing.T) {
	grandparent := NewEnv(Bindings{"grandparent": envTestValue("grandparent")})
	parent := Env{
		inner: Bindings{"parent": envTestValue("parent"), "shadowed": envTestValue("parent")},
		outer: &grandparent,
	}
	env := Env{
		inner: Bindings{"local": envTestValue("local"), "shadowed": envTestValue("local")},
		outer: &parent,
	}

	for key, want := range map[string]Value{
		"local":       envTestValue("local"),
		"parent":      envTestValue("parent"),
		"grandparent": envTestValue("grandparent"),
		"shadowed":    envTestValue("local"),
	} {
		if got := env.Get(key); got != want {
			t.Errorf("Get(%q) = %v, want %v", key, got, want)
		}
	}
	if got := env.Get("missing"); got != nil {
		t.Errorf("Get(%q) = %v, want nil", "missing", got)
	}
}

func TestEnvSetUpdatesOnlyCurrentScope(t *testing.T) {
	parent := NewEnv(Bindings{
		"shared": envTestValue("parent"),
		"parent": envTestValue("parent"),
	})
	env := Env{inner: Bindings{}, outer: &parent}

	env.Set("shared", envTestValue("child"))
	env.Set("local", envTestValue("local"))

	if got := env.Get("shared"); got != envTestValue("child") {
		t.Errorf("Get(%q) = %v, want %v", "shared", got, envTestValue("child"))
	}
	if got := parent.Get("shared"); got != envTestValue("parent") {
		t.Errorf("parent Get(%q) = %v, want %v", "shared", got, envTestValue("parent"))
	}
	if got := env.inner["local"]; got != envTestValue("local") {
		t.Errorf("inner[%q] = %v, want %v", "local", got, envTestValue("local"))
	}
}

func TestBindingsAsEnvSharesBindings(t *testing.T) {
	bindings := Bindings{"key": envTestValue("before")}
	env := bindings.AsEnv()

	bindings["key"] = envTestValue("after")
	if got := env.Get("key"); got != envTestValue("after") {
		t.Errorf("Get(%q) after changing bindings = %v, want %v", "key", got, envTestValue("after"))
	}

	env.Set("key", envTestValue("set"))
	if got := bindings["key"]; got != envTestValue("set") {
		t.Errorf("bindings[%q] after Set = %v, want %v", "key", got, envTestValue("set"))
	}
}
