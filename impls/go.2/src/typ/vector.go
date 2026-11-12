package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/fun"
	"github.com/mooss/mal/src/ut"
)

// VectorForm is an unevaluated sequence of values.
type VectorForm []core.Value

func (v VectorForm) Equal(other core.Value) bool { return sliceEq(v, other, core.Value.Equal) }

func (v VectorForm) Eval(envr core.Env) core.Value {
	res := make(Vector, len(v))
	for i, el := range v {
		value := Resolve(envr, el)
		if ut.Is[Fatal](value) {
			return value
		}
		res[i] = value
	}
	return res
}

func (v VectorForm) String() string {
	return "[" + strings.Join(fun.Map(core.Value.String, v), " ") + "]"
}

func (VectorForm) Type() core.TypeName { return core.VECTOR_FORM }

// Vector is a sequence of values.
type Vector []core.Value

func (v Vector) Equal(other core.Value) bool { return sliceEq(v, other, core.Value.Equal) }
func (v Vector) Eval(core.Env) core.Value    { return v }

func (v Vector) String() string {
	return "[" + strings.Join(fun.Map(core.Value.String, v), " ") + "]"
}

func (Vector) Type() core.TypeName { return core.VECTOR }
