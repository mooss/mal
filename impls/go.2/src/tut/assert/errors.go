package assert

import (
	"strings"
	"testing"
)

// Err fails if err == nil.
func Err(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ErrContains fails if err is nil or does not contain sub.
// An empty sub only checks for an error.
func ErrContains(t *testing.T, err error, sub string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", sub)
	}
	if sub != "" && !strings.Contains(err.Error(), sub) {
		t.Fatalf("%q does not contain %q", err.Error(), sub)
	}
}

// Must returns v when err is nil, otherwise fails.
func Must[T any](t *testing.T, v T, err error) T {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return v
}

// NoErr fails if err != nil.
func NoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
