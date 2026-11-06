package core

import "strings"

// List is a sequence of values.
type List []Type

func (List) Type() TypeName   { return LIST }
func (l List) String() string { return "(" + strings.Join(_map(Type.String, l), " ") + ")" }
