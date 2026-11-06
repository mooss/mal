package core

import "strings"

// Vector is a sequence of values in brackets.
type Vector []Type

func (Vector) Type() TypeName   { return VECTOR }
func (v Vector) String() string { return "[" + strings.Join(_map(Type.String, v), " ") + "]" }
