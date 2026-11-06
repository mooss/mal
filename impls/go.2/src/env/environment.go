// Package env implements environments.
package env

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

type _fun = typ.Function
type _int = typ.Integer

var Initial core.Env = map[string]core.Type{
	"+": _fun(func(l _int, r _int) _int { return l + r }),
	"-": _fun(func(l _int, r _int) _int { return l - r }),
	"*": _fun(func(l _int, r _int) _int { return l * r }),
	"/": _fun(func(l _int, r _int) _int { return l / r }), // Will panic on r == 0.
}
