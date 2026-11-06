package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
)

// List is a sequence of values.
type List []core.Type

func (List) Type() core.TypeName { return core.LIST }
func (lst List) String() string {
	return "(" + strings.Join(_map(core.Type.String, lst), " ") + ")"
}

func (lst List) Eval(envr core.Env) core.Type {
	if len(lst) == 0 {
		return lst
	}
	if len(lst) != 3 {
		msg := "only list of 3 symbols can be evaluated right now (got %d)"
		return Errorf(msg, len(lst))
	}

	// Apply phase.
	fun, cerr := as[_fun](lst[0].Eval(envr))
	if cerr != nil {
		return *cerr
	}

	arg1, cerr := as[_int](lst[1].Eval(envr))
	if cerr != nil {
		return *cerr
	}

	arg2, cerr := as[_int](lst[2].Eval(envr))
	if cerr != nil {
		return *cerr
	}

	return fun(arg1, arg2)
}
