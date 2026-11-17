package assert

import (
	"reflect"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/tut/eq"
)

// Compare fails when compare(got, want) returns false.
func Compare[T any](t *testing.T, got, want T, name string, compare func(a, b T) bool) {
	t.Helper()
	if !compare(got, want) {
		t.Fatalf("%s(%#v, %#v) failed", name, got, want)
	}
}
func (tr T) Compare[T any](got, want T, name string, compare func(a, b T) bool) {
	tr.T.Helper()
	Compare(tr.T, got, want, name, compare)
}

// Equal fails if got != want according to reflect.DeepEqual.
func Equal[T any](t *testing.T, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Equal failed:\ngot  %#v\nwant %#v", got, want)
	}
}
func (tr T) Equal[T any](got, want T) {
	tr.T.Helper()
	Equal(tr.T, got, want)
}

// MalDifferent fails unless MalDifferent(got, want).
func MalDifferent(t *testing.T, got, want core.Value) {
	t.Helper()
	Compare(t, got, want, "eq.MalNot", eq.MalNot)
}
func (tr T) MalDifferent(got, want core.Value) {
	tr.T.Helper()
	MalDifferent(tr.T, got, want)
}

// MalEqual fails unless MalEqual(got, want).
func MalEqual(t *testing.T, got, want core.Value) {
	t.Helper()
	Compare(t, got, want, "eq.Mal", eq.Mal)
}
func (tr T) MalEqual(got, want core.Value) {
	tr.T.Helper()
	MalEqual(tr.T, got, want)
}

// MalEqualSlice fails unless MalEqualSlices(got, want).
func MalEqualSlice(t *testing.T, got, want []core.Value) {
	t.Helper()
	Compare(t, got, want, "eq.MalSlice", eq.MalSlice)
}
func (tr T) MalEqualSlice(got, want []core.Value) {
	tr.T.Helper()
	MalEqualSlice(tr.T, got, want)
}
