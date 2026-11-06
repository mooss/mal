package core

import "strings"

type Vector []Type

func (Vector) Type() TypeName   { return VECTOR }
func (v Vector) String() string { return "[" + strings.Join(_map(Type.String, v), " ") + "]" }
