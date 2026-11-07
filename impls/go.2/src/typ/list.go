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
	fun, err := as[Function](lst[0].Eval(envr))
	if err.NotNil() {
		return err
	}

	args := make([]core.Value, len(lst)-1)
	for i, arg := range lst[1:] {
		arg = arg.Eval(envr)
		if ut.Is[Fatal](arg) {
			return arg
		}
		args[i] = arg
	}

	return fun(args...)
}
