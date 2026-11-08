package typ

import (
	"github.com/mooss/mal/src/core"
)

// Function represents a plain function.
// When encountered as the first elements of a list, the rest of the list will be evaluated and
// those evaluated arguments will be passed to the function.
type Function func(...core.Value) core.Value

func (Function) Type() core.TypeName          { return core.FUNCTION }
func (Function) String() string               { return "<internal function>" }
func (fun Function) Eval(core.Env) core.Value { return fun }

////////////////////////
// Conversion helpers //

func II[A, B, R core.Value](fun func(A, B) R) Function {
	return func(args ...core.Value) core.Value {
		a, b, err := as2Args[A, B](args)
		if err.NotNil() {
			return err
		}

		return fun(a, b)
	}
}
