package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

// isEmpty returns true when the given container has no elements.
func isEmpty(arg core.Value) core.Value {
	size, err := length(arg)
	if err.NotNil() {
		return err
	}
	return typ.Bool(size == 0)
}

// isFalse returns true when arg is a boolean and its value is false.
func isFalse(arg core.Value) typ.Bool {
	boo, isBool := arg.(typ.Bool)
	if isBool {
		return !boo
	}
	return false
}

// isSequential returns true when arg is a list or vector.
func isSequential(arg core.Value) typ.Bool {
	switch arg.(type) {
	case typ.Vector, typ.List:
		return true
	default:
		return false
	}
}

// isTrue returns true when arg is a boolean and its value is true.
func isTrue(arg core.Value) typ.Bool {
	boo, isBool := arg.(typ.Bool)
	if isBool {
		return boo
	}
	return false
}

// isX returns true when the type of arg is X.
func isX[X core.Value](arg core.Value) typ.Bool { return typ.Bool(ut.Is[X](arg)) }
