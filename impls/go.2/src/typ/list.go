package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
)

// List is a sequence of values.
type List []core.Type

func (List) Type() core.TypeName { return core.LIST }
func (l List) String() string    { return "(" + strings.Join(_map(core.Type.String, l), " ") + ")" }
