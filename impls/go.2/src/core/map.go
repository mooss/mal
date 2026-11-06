package core

import "strings"

// Map is a mapping from values to values.
type Map map[Type]Type

func (Map) Type() TypeName   { return MAP }
func (m Map) String() string { return "{" + strings.Join(_mapflat(Type.String, m), " ") + "}" }

// MapLit is an unevaluated sequence of map keys and values.
type MapLit []Type

func (MapLit) Type() TypeName   { return MAP_LITERAL }
func (m MapLit) String() string { return "{" + strings.Join(_map(Type.String, m), " ") + "}" }
