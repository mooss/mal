package typ_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func TestAsSeqEmptyNoError(t *testing.T) {
	seq, getErr := typ.AsSeq[typ.Integer]()
	got := slices.Collect(seq)
	if len(got) != 0 {
		t.Fatalf("collected %d values, want 0", len(got))
	}
	if err := getErr(); err.NotNil() {
		t.Fatalf("getter = %q, want no error for empty input", err.String())
	}
	// Iterating again stays empty with no error.
	got2 := slices.Collect(seq)
	if len(got2) != 0 {
		t.Fatalf("second iteration collected %d values, want 0", len(got2))
	}
	if err := getErr(); err.NotNil() {
		t.Fatalf("getter after second iteration = %q, want no error", err.String())
	}
}

func TestAsSeqAllValidNoError(t *testing.T) {
	seq, getErr := typ.AsSeq[typ.Integer](typ.Integer(1), typ.Integer(2), typ.Integer(3))
	got := slices.Collect(seq)
	if len(got) != 3 {
		t.Fatalf("collected %d values, want 3", len(got))
	}
	for i, want := range []typ.Integer{typ.Integer(1), typ.Integer(2), typ.Integer(3)} {
		if got[i] != want {
			t.Errorf("got[%d] = %s, want %s", i, got[i].String(), want.String())
		}
	}
	if err := getErr(); err.NotNil() {
		t.Fatalf("getter = %q, want no error for all valid inputs", err.String())
	}
}

func TestAsSeqBadTypeStopsIteration(t *testing.T) {
	seq, getErr := typ.AsSeq[typ.Integer](
		typ.Integer(1),
		typ.String("oops"),
		typ.Integer(3),
	)
	got := slices.Collect(seq)
	if len(got) != 1 || got[0] != typ.Integer(1) {
		t.Fatalf("collected = %v, want [1] (stop at first bad type)", got)
	}
	err := getErr()
	if !err.NotNil() {
		t.Fatal("getter returned no error, want Fatal for bad type")
	}
	if !strings.Contains(strings.ToLower(err.String()), "expected") {
		t.Fatalf("getter = %q, want message containing 'expected'", err.String())
	}
	// Iteration after fatal stays stopped.
	for range seq {
		t.Fatal("iterator yielded after fatal, want stopped")
	}
	if err2 := getErr(); !err2.NotNil() {
		t.Fatal("getter after second iteration returned no error, want sticky Fatal")
	}
}

func TestAsSeqFatalForwardedAsIs(t *testing.T) {
	fatal := typ.Fatalf("boom %d", 42)
	seq, getErr := typ.AsSeq[typ.Integer](
		typ.Integer(1),
		fatal,
		typ.Integer(3),
	)
	got := slices.Collect(seq)
	if len(got) != 1 || got[0] != typ.Integer(1) {
		t.Fatalf("collected = %v, want [1] (stop at Fatal input)", got)
	}
	err := getErr()
	if !err.NotNil() {
		t.Fatal("getter returned no error, want forwarded Fatal")
	}
	if err.String() != fatal.String() {
		t.Fatalf("getter = %q, want forwarded Fatal %q as is", err.String(), fatal.String())
	}
}

func TestAsSeqResumeAfterPartialConsume(t *testing.T) {
	seq, getErr := typ.AsSeq[typ.Integer](
		typ.Integer(1),
		typ.Integer(2),
		typ.Integer(3),
	)
	first, ok := ut.First(seq)
	if !ok {
		t.Fatal("First found no value, want 1")
	}
	if first != typ.Integer(1) {
		t.Fatalf("First = %s, want 1", first.String())
	}
	if err := getErr(); err.NotNil() {
		t.Fatalf("getter after First = %q, want no error", err.String())
	}
	// Reduce the remainder; must resume at 2, not restart at 1.
	sum := ut.Reduce(func(acc typ.Integer, v typ.Integer) typ.Integer {
		return acc + v
	}, 0, seq)
	if sum != 5 {
		t.Fatalf("Reduce after First = %d, want 5 (2+3), iterator did not resume)", sum)
	}
	if err := getErr(); err.NotNil() {
		t.Fatalf("getter after Reduce = %q, want no error", err.String())
	}
	// Fully consumed: further iteration yields nothing.
	if _, ok := ut.First(seq); ok {
		t.Fatal("iterator yielded after full consume, want exhausted")
	}
}
