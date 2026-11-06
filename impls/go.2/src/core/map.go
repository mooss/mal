package core

import "strings"

type Map map[Type]Type

func (Map) Type() TypeName   { return MAP }
func (m Map) String() string { return "{" + strings.Join(_mapflat(Type.String, m), " ") + "}" }

type MapLit []Type

func (MapLit) Type() TypeName   { return MAP_LITERAL }
func (m MapLit) String() string { return "{" + strings.Join(_map(Type.String, m), " ") + "}" }
