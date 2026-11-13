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

// Equal fails if got != want according to reflect.DeepEqual.
func Equal[T any](t *testing.T, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// MalDifferent fails unless MalDifferent(got, want).
func MalDifferent(t *testing.T, got, want core.Value) {
	t.Helper()
	Compare(t, got, want, "eq.MalNot", eq.MalNot)
}

// MalEqual fails unless MalEqual(got, want).
func MalEqual(t *testing.T, got, want core.Value) {
	t.Helper()
	Compare(t, got, want, "eq.Mal", eq.Mal)
}

// MalEqualSlice fails unless MalEqualSlices(got, want).
func MalEqualSlice(t *testing.T, got, want []core.Value) {
	t.Helper()
	Compare(t, got, want, "eq.MalSlice", eq.MalSlice)
}
