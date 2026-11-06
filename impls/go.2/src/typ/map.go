package typ

import (
	"strings"

	"github.com/mooss/mal/src/core"
)

// Map is a mapping from values to values.
type Map map[core.Type]core.Type

func (Map) Type() core.TypeName { return core.MAP }
func (m Map) String() string {
	return "{" + strings.Join(_mapflat(core.Type.String, m), " ") + "}"
}

// MapLit is an unevaluated sequence of map keys and values.
type MapLit []core.Type

func (MapLit) Type() core.TypeName { return core.MAP_LITERAL }
func (m MapLit) String() string {
	return "{" + strings.Join(_map(core.Type.String, m), " ") + "}"
}
