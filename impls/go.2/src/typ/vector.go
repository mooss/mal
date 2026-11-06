package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
)

// Vector is a sequence of values in brackets.
type Vector []core.Type

func (Vector) Type() core.TypeName { return core.VECTOR }
func (v Vector) String() string {
	return "[" + strings.Join(_map(core.Type.String, v), " ") + "]"
}
