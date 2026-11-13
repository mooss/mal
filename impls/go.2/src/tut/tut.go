// Package tut provides test utilities.
package tut

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/typ"
)

// MustEval parses one form and resolves it, panicking on parse error.
func MustEval(envr core.Env, input string) core.Value {
	program, err := reader.ParseOne(input)
	if err != nil {
		panic(err)
	}

	return typ.Resolve(envr, program)
}

// MustEvalAs evaluates input and asserts the result has type T.
func MustEvalAs[T core.Value](t *testing.T, env core.Env, input string) T {
	t.Helper()
	got := MustEval(env, input)
	v, ok := got.(T)
	if !ok {
		var zero T
		t.Fatalf("MustEval(%q) type = %T, want %T", input, got, zero)
	}
	return v
}

// Case is a string-input table entry with expected value and error flag.
type Case[W any] struct {
	Name    string
	Input   string
	Want    W
	WantErr bool
}

// Run executes each case as a subtest.
func Run[W any](t *testing.T, cases []Case[W], fn func(*testing.T, Case[W])) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Helper()
			fn(t, c)
		})
	}
}

// SimpleCase is a success-only table entry without error flag.
type SimpleCase[W any] struct {
	Name string
	Prog string
	Want W
}

// RunSimple executes each simple case as a subtest.
func RunSimple[W any](t *testing.T, cases []SimpleCase[W], fn func(*testing.T, SimpleCase[W])) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Helper()
			fn(t, c)
		})
	}
}

// MalEqual reports Mal value equality with .Equal.
// Will panic on nil values, this is intentional, Mal should never return nil values.
func MalEqual(a, b core.Value) bool {
	if a == nil || b == nil {
		panic(fmt.Sprintf("MalEqual(%#v, %#v) got nil value", a, b))
	}
	return a.Equal(b)
}

// MalDifferent reports Mal value inequality with .Equal.
// Will panic on nil values, this is intentional, Mal should never return nil values.
func MalDifferent(a, b core.Value) bool {
	if a == nil || b == nil {
		panic(fmt.Sprintf("MalDifferent(%#v, %#v) got nil value", a, b))
	}
	return !a.Equal(b)
}

// EqualSlicesWith reports element-wise equality using eq.
// It is generic over any element type and makes EqualSlices reusable.
func EqualSlicesWith[T any](a, b []T, eq func(a, b T) bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !eq(a[i], b[i]) {
			return false
		}
	}
	return true
}

// MalEqualSlices reports element-wise MalEqual.
func MalEqualSlices(a, b []core.Value) bool {
	return EqualSlicesWith(a, b, MalEqual)
}

// Capture runs fn with stdout piped and returns what was printed.
func Capture(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("pipe writer Close() error = %v", err)
	}
	os.Stdout = old
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("pipe reader Close() error = %v", err)
	}
	return buf.String()
}
