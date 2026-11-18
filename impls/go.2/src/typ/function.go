package typ

import (
	"github.com/mooss/mal/src/core"
)

// Function is a function (primitive or defined in MAL itself).
// When encountered as the first elements of a list, the rest of the list will be evaluated and
// those evaluated arguments will be passed to the function.
type Function func(...core.Value) core.Value

func (Function) Equal(core.Value) bool        { return false }
func (fun Function) Eval(core.Env) core.Value { return fun }
func (Function) String() string               { return "<function>" }
func (Function) Type() core.TypeName          { return core.FUNCTION }

// Macro is a function manipulating the AST.
type Macro func(...core.Value) core.Value

func (Macro) Equal(core.Value) bool        { return false }
func (mac Macro) Eval(core.Env) core.Value { return mac }
func (Macro) String() string               { return "<macro>" }
func (Macro) Type() core.TypeName          { return core.MACRO }

////////////////////////
// Conversion helpers //

func I[A, R core.Value](fun func(A) R) Function {
	return func(args ...core.Value) core.Value {
		a, err := as1Arg[A](args)
		if err.NotNil() {
			return err
		}

		return fun(a)
	}
}

func II[A, B, R core.Value](fun func(A, B) R) Function {
	return func(args ...core.Value) core.Value {
		a, b, err := as2Args[A, B](args)
		if err.NotNil() {
			return err
		}

		return fun(a, b)
	}
}
