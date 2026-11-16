package assert

import (
	"strings"
	"testing"
)

// ErrorContains fails if err is nil or does not contain sub.
// An empty sub only checks for an error.
func ErrorContains(t *testing.T, err error, sub string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", sub)
	}
	if sub != "" && !strings.Contains(err.Error(), sub) {
		t.Fatalf("%q does not contain %q", err.Error(), sub)
	}
}
func (tr T) ErrorContains(err error, sub string) {
	tr.T.Helper()
	ErrorContains(tr.T, err, sub)
}

// Must returns v when err is nil, otherwise fails.
func Must[T any](t *testing.T, v T, err error) T {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return v
}
func (tr T) Must[T any](v T, err error) T {
	tr.T.Helper()
	return Must(tr.T, v, err)
}

// MustError fails if err == nil.
func MustError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
func (tr T) MustError(err error) {
	tr.T.Helper()
	MustError(tr.T, err)
}

// NoError fails if err != nil.
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func (tr T) NoError(err error) {
	tr.T.Helper()
	NoError(tr.T, err)
}
