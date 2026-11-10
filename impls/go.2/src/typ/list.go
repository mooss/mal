package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/fun"
)

// List is a sequence of values.
type List []core.Value

func (List) Type() core.TypeName { return core.LIST }
func (lst List) String() string {
	return "(" + strings.Join(fun.Map(core.Value.String, lst), " ") + ")"
}

func (lst List) Eval(envr core.Env) core.Value {
	if len(lst) == 0 {
		return lst
	}

	// Apply phase.
	op := lst[0].Eval(envr)
	switch op := op.(type) {
	case Function:
		return evalFun(envr, op, lst[1:])
	case Special:
		return op(envr, lst[1:]...)
	}

	return Expected("callable", op)
}

////////////////////////////
// Eval-related utilities //

// evalFun evaluates a function by calling it with its evaluated arguments.
// Any fatal argument encountered stops the evaluation and is returned as-is.
func evalFun(envr core.Env, fn Function, rest []core.Value) core.Value {
	args, err := EvalArgs(envr, rest...)
	if err.NotNil() {
		return err
	}

	return fn(args...)
}
