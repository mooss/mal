package assert

import (
	"strings"
	"testing"

	"github.com/mooss/mal/src/core"
)

// StringContains fails if got does not contain the given substring.
func StringContains(t *testing.T, got, substring string) {
	t.Helper()
	if !strings.Contains(got, substring) {
		t.Fatalf("%q does not contain %q", got, substring)
	}
}

// StringRepr fails if v.String() != want.
func StringRepr(t *testing.T, v core.Value, want string) {
	t.Helper()
	if v == nil {
		t.Fatalf("got nil, want %q", want)
	}
	if got := v.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
