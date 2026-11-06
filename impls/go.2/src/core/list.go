package core

import "strings"

type List []Type

func (List) Type() TypeName   { return LIST }
func (l List) String() string { return "(" + strings.Join(_map(Type.String, l), " ") + ")" }
