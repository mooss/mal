package envs

import (
	"maps"
	"slices"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/fun"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

// assoc copies the first argument as a map and adds new association from the rest of the arguments.
func assoc(args ...core.Value) core.Value {
	if err := typ.ExpectAtLeastNArgs(1, args); err.NotNil() {
		return err
	}

	source, err := typ.As[typ.Map](args[0])
	if err.NotNil() {
		return err
	}

	if len(args)%2 != 1 { // != 1 to take the first argument into account.
		return OddMapArgsErr
	}

	return extendMap(maps.Clone(source), args[1:]...)
}

// dissoc (map, ...keys) returns a copy of map that does not contains any of ...keys.
func dissoc(args ...core.Value) core.Value {
	if err := typ.ExpectAtLeastNArgs(1, args); err.NotNil() {
		return err
	}
	source, err := typ.As[typ.Map](args[0])
	if err.NotNil() {
		return err
	}

	keys := ut.NewSetI(fun.Filter(typ.ValidMapKey, slices.Values(args[1:])))
	// Allocate the strict minimum to try and avoid repeated growth.
	res := make(typ.Map, max(0, len(source)-len(keys)))
	for key, value := range source {
		if !keys.Contains(key) {
			res[key] = value
		}
	}
	return res
}

// get (map|nil, key) returns the value associated to key in map, nil if it does not exist.
func get(dict, key core.Value) core.Value {
	switch typed := dict.(type) {
	case typ.Nil:
		return typ.Nil{}
	case typ.Map:
		return typed.MalGet(key)
	default:
		return typ.Expected("nil or map", dict)
	}
}

// mkMap creates a new map from even key, value... arguments.
func mkMap(args ...core.Value) core.Value {
	if len(args)%2 != 0 {
		return OddMapArgsErr
	}

	return extendMap(make(typ.Map, len(args)/2), args...)
}
