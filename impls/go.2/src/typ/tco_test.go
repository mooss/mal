package typ_test

import (
	"runtime"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/typ"
)

func tailCallEnv() (core.Env, *int) {
	maxFrames := new(int)
	return envs.MalSpecial.Nest(core.Bindings{
		"zero?": typ.Function(func(args ...core.Value) core.Value {
			pcs := make([]uintptr, 256)
			if frames := runtime.Callers(0, pcs); frames > *maxFrames {
				*maxFrames = frames
			}
			return typ.Bool(args[0].(typ.Integer) == 0)
		}),
		"nonzero?": typ.Function(func(args ...core.Value) core.Value {
			pcs := make([]uintptr, 256)
			if frames := runtime.Callers(0, pcs); frames > *maxFrames {
				*maxFrames = frames
			}
			return typ.Bool(args[0].(typ.Integer) != 0)
		}),
		"dec": typ.Function(func(args ...core.Value) core.Value {
			return typ.Integer(args[0].(typ.Integer) - 1)
		}),
	}), maxFrames
}

// All those tests fail when replacing typ.Tail with a direct evaluation.
// When adding new ones, make sure those failures happen on the new tests too.
func TestTailCallsDoNotGrowStack(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "if false branch",
			input: `(do
				(def! loop (fn* (n)
					(if (zero? n) 0 (loop (dec n)))))
				(loop 50000))`,
		},
		{
			name: "if true branch",
			input: `(do
				(def! loop (fn* (n)
					(if (nonzero? n) (loop (dec n)) 0)))
				(loop 50000))`,
		},
		{
			name: "do",
			input: `(do
				(def! loop (fn* (n)
					(if (zero? n) 0 (do nil (loop (dec n))))))
				(loop 50000))`,
		},
		{
			name: "let*",
			input: `(do
				(def! loop (fn* (n)
					(if (zero? n) 0
						(let* [next (dec n)] (loop next)))))
				(loop 50000))`,
		},
		{
			name: "variadic closure",
			input: `(do
				(def! loop (fn* (n & rest)
					(if (zero? n) 0 (loop (dec n) nil))))
				(loop 50000))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, maxFrames := tailCallEnv()
			got := tut.MustEval(env, tt.input)
			if got != core.Value(typ.Integer(0)) {
				t.Fatalf("MustEval() = %v, want 0", got)
			}
			if *maxFrames > 100 {
				t.Fatalf("stack grew to %d frames during tail recursion", *maxFrames)
			}
		})
	}
}
