package typ

import (
	"fmt"
	"os"
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/fun"
)

// canDebug guards the ability to enable debug mode.
const canDebug = false

// List is a sequence of values.
type List []core.Value

func NewList(elements ...core.Value) List { return List(elements) }

func (lst List) Equal(other core.Value) bool { return sliceEq(lst, other, core.Value.Equal) }

func (lst List) Eval(envr core.Env) core.Value {
	if canDebug && IsTrue(envr.Get("DEBUG-MODE")) {
		fmt.Fprintln(os.Stderr, ":APPLY", lst)
	}

	if len(lst) == 0 {
		return lst
	}

	op := Resolve(envr, lst[0])
	switch op := op.(type) {
	case Function:
		return evalFun(envr, op, lst[1:])
	case Special:
		return op(envr, lst[1:]...)
	}

	return Expected("callable", op)
}

func (lst List) String() string {
	return "(" + strings.Join(fun.Map(core.Value.String, lst), " ") + ")"
}

func (List) Type() core.TypeName { return core.LIST }

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
