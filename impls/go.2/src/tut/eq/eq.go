package eq

import (
	"fmt"
	"slices"

	"github.com/mooss/mal/src/core"
)

// Mal reports Mal value equality with .Equal.
// Will panic on nil values, this is intentional, Mal should never return nil values.
func Mal(a, b core.Value) bool {
	if a == nil || b == nil {
		panic(fmt.Sprintf("Mal(%#v, %#v) got nil value", a, b))
	}
	return a.Equal(b)
}

// MalNot reports Mal value inequality with .Equal.
// Will panic on nil values, this is intentional, Mal should never return nil values.
func MalNot(a, b core.Value) bool {
	if a == nil || b == nil {
		panic(fmt.Sprintf("MalNot(%#v, %#v) got nil value", a, b))
	}
	return !a.Equal(b)
}

// MalSlice reports element-wise Mal equality.
func MalSlice(a, b []core.Value) bool {
	return slices.EqualFunc(a, b, Mal)
}
