package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

// apply calls the first argument as a function, with the rest of the arguments forwarded to it.
// Its last argument must be a list or vector, to which the in-between arguments will be prepended.
func apply(args ...core.Value) core.Value {
	if err := typ.ExpectAtLeastNArgs(2, args); err.NotNil() {
		return err
	}

	fn, err := typ.AsFunOrMac(args[0])
	if err.NotNil() {
		return err
	}

	last, err := typ.AsListOrVec(args[len(args)-1])
	if err.NotNil() {
		return err
	}

	params := make(typ.List, len(last)+len(args)-2)
	copy(params, args[1:len(args)-1])
	copy(params[len(args)-2:], last)
	return fn(params...)
}

// mapfn calls its first argument on all the elements of its second arguments and returns the result
// as a new list.
func mapfn(fn core.Value, list core.Value) core.Value {
	rawfun, err := typ.AsFunOrMac(fn)
	if err.NotNil() {
		return err
	}

	golist, err := typ.AsListOrVec(list)
	if err.NotNil() {
		return err
	}

	res := make(typ.List, len(golist))
	for i, el := range golist {
		elmapped := typ.EnsureResolved(rawfun(el))
		if ut.Is[typ.Fatal](elmapped) {
			return elmapped
		}
		res[i] = elmapped
	}

	return res
}
