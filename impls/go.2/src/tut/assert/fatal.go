package assert

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

// Fatal fails unless v is Fatal and returns it.
func Fatal(t *testing.T, v core.Value) typ.Fatal {
	t.Helper()
	if f, ok := v.(typ.Fatal); ok {
		return f
	}
	if _, ok := v.(*typ.Fatal); ok {
		panic("got fatal pointer")
	}

	//revive:disable:unreachable-code (>< I still have to return don't I?)
	t.Fatalf("got %v (%T), want Fatal", v, v)
	var zero typ.Fatal
	return zero
}
func (tr T) Fatal(v core.Value) typ.Fatal {
	tr.T.Helper()
	return Fatal(tr.T, v)
}

// Ok fails if v is Fatal and returns v.
func Ok(t *testing.T, v core.Value) core.Value {
	t.Helper()
	if IsFatal(v) {
		t.Fatalf("got Fatal %v, want success", v)
	}
	return v
}
func (tr T) Ok(v core.Value) core.Value {
	tr.T.Helper()
	return Ok(tr.T, v)
}
