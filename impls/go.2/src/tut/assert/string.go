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
func (tr T) StringContains(got, substring string) {
	tr.T.Helper()
	StringContains(tr.T, got, substring)
}

// StringEndsWith fails if got does not end with suffix.
func StringEndsWith(t *testing.T, got, suffix string) {
	t.Helper()
	if !strings.HasSuffix(got, suffix) {
		t.Fatalf("%q does not end with %q", got, suffix)
	}
}
func (tr T) StringEndsWith(got, suffix string) {
	tr.T.Helper()
	StringEndsWith(tr.T, got, suffix)
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
func (tr T) StringRepr(v core.Value, want string) {
	tr.T.Helper()
	StringRepr(tr.T, v, want)
}

// StringStartsWith fails if got does not start with prefix.
func StringStartsWith(t *testing.T, got, prefix string) {
	t.Helper()
	if !strings.HasPrefix(got, prefix) {
		t.Fatalf("%q does not start with %q", got, prefix)
	}
}
func (tr T) StringStartsWith(got, prefix string) {
	tr.T.Helper()
	StringStartsWith(tr.T, got, prefix)
}
