package assert

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

// Is asserts v has type T and returns it.
func Is[T core.Value](t *testing.T, v core.Value) T {
	t.Helper()
	got, ok := v.(T)
	if !ok {
		var zero T
		t.Fatalf("type = %T, want %T", v, zero)
	}
	return got
}

// IsFatal reports whether v is a Fatal value.
func IsFatal(v core.Value) bool {
	if v == nil {
		return false
	}
	if _, ok := v.(typ.Fatal); ok {
		return true
	}
	if _, ok := v.(*typ.Fatal); ok {
		panic("got fatal pointer")
	}
	return false
}
