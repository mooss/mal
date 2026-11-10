// Package envs provides pre-defined environments.
package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

type _int = typ.Integer

var MalDefault = core.CatEnv(MalSpecial, Math)
