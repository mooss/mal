package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
)

// Map is a mapping from values to values.
type Map map[core.Value]core.Value

func (Map) Type() core.TypeName { return core.MAP }
func (m Map) String() string {
	return "{" + strings.Join(_mapflat(core.Value.String, m), " ") + "}"
}

// TBD
func (m Map) Eval(core.Env) core.Value { return m }

// MapLit is an unevaluated sequence of map keys and values.
type MapLit []core.Value

func (MapLit) Type() core.TypeName { return core.MAP_LITERAL }
func (m MapLit) String() string {
	return "{" + strings.Join(_map(core.Value.String, m), " ") + "}"
}

// TBD
func (m MapLit) Eval(core.Env) core.Value { return m }
