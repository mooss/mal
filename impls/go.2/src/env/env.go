// Package env implements environments.
package env

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

type _int = typ.Integer

// bint transforms a binary int function into a typ.Function
var bint = typ.II[_int, _int, _int]

var Initial = core.Bindings{
	"+": bint(func(l _int, r _int) _int { return l + r }),
	"-": bint(func(l _int, r _int) _int { return l - r }),
	"*": bint(func(l _int, r _int) _int { return l * r }),
	"/": bint(func(l _int, r _int) _int { return l / r }), // Will panic on r == 0.
}.AsEnv()
