package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/ut"
)

// Vector is a sequence of values in brackets.
type Vector []core.Type

func (Vector) Type() core.TypeName { return core.VECTOR }
func (v Vector) String() string {
	return "[" + strings.Join(_map(core.Type.String, v), " ") + "]"
}

func (v Vector) Eval(envr core.Env) core.Type {
	res := make(_vec, len(v))
	for i, el := range v {
		value := el.Eval(envr)
		if ut.Is[Error](value) {
			return value
		}
		res[i] = value
	}
	return res
}
