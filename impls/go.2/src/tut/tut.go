// Package tut provides test utilities.
package tut

import (
	"bytes"
	"io"
	"os"
	"reflect"
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
	Name  string
	Input string
	Want  W
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

// Equal reports strict Mal value equality: same Go type and Equal.
// Main difference should be that vector and lists cannot be equal anymore.
func Equal(a, b core.Value) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	return a.Equal(b)
}

// EqualSlices reports element-wise Equal.
func EqualSlices(a, b []core.Value) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !Equal(a[i], b[i]) {
			return false
		}
	}
	return true
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
