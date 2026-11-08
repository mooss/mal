package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/ut"
)

// List is a sequence of values.
type List []core.Value

func (List) Type() core.TypeName { return core.LIST }
func (lst List) String() string {
	return "(" + strings.Join(_map(core.Value.String, lst), " ") + ")"
}

func (lst List) Eval(envr core.Env) core.Value {
	if len(lst) == 0 {
		return lst
	}
	if len(lst) != 3 {
		msg := "only list of 3 symbols can be evaluated right now (got %d)"
		return Fatalf(msg, len(lst))
	}

	// Apply phase.
	op := lst[0].Eval(envr)
	switch op := op.(type) {
	case Function:
		return evalFun(envr, op, lst[1:])
	case Special:
		return op(envr, lst[1:]...)
	case Fatal:
		return op
	}

	return expected("callable", op)
}

////////////////////////////
// Eval-related utilities //

// evalFun evaluates a function by calling it with its evaluated arguments.
// Any fatal argument encountered stops the evaluation and is returned as-is.
func evalFun(envr core.Env, fun Function, rest []core.Value) core.Value {
	args := make([]core.Value, len(rest))
	for i, arg := range rest {
		arg = arg.Eval(envr)
		if ut.Is[Fatal](arg) {
			return arg
		}
		args[i] = arg
	}

	return fun(args...)
}
