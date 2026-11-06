package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
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
	fun, cerr := as[Function](lst[0].Eval(envr))
	if cerr != nil {
		return *cerr
	}

	arg1, cerr := as[Integer](lst[1].Eval(envr))
	if cerr != nil {
		return *cerr
	}

	arg2, cerr := as[Integer](lst[2].Eval(envr))
	if cerr != nil {
		return *cerr
	}

	return fun(arg1, arg2)
}
