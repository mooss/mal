package assert

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/tut/eq"
)

// EnvHas fails unless env.EnvHas(key) strictly equals want.
func EnvHas(t *testing.T, env core.Env, key string, want core.Value) {
	t.Helper()
	if got := env.Get(key); !eq.Mal(got, want) {
		t.Fatalf("Get(%q) = %v, want %v", key, got, want)
	}
}

// EnvLacks fails if any key is bound.
func EnvLacks(t *testing.T, env core.Env, keys ...string) {
	t.Helper()
	for _, k := range keys {
		if got := env.Get(k); got != nil {
			t.Fatalf("Get(%q) = %v, want nil", k, got)
		}
	}
}
