package typ

import (
	"fmt"

	"github.com/mooss/mal/src/core"
)

type Function func(...core.Value) core.Value

func (Function) Type() core.TypeName          { return core.FUNCTION }
func (fun Function) String() string           { return fmt.Sprintf("<internal %T>", fun) }
func (fun Function) Eval(core.Env) core.Value { return fun }

////////////////////////
// Conversion helpers //

func II[A, B, R core.Value](fun func(A, B) R) Function {
	return func(args ...core.Value) core.Value {
		if len(args) != 2 {
			return Fatalf("expected 2 args, got %d", len(args))
		}

		a, err := as[A](args[0])
		if err.NotNil() {
			return err
		}
		b, err := as[B](args[1])
		if err.NotNil() {
			return err
		}

		return fun(a, b)
	}
}
