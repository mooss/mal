package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

// deref returns the value currently held by the atom.
func deref(atom *typ.Atom) core.Value { return atom.Ref }

// mkAtom creates a new atom from a value.
func mkAtom(arg core.Value) *typ.Atom { return &typ.Atom{Ref: arg} }

// reset replaces the value held by the atom and returns it.
func reset(atom *typ.Atom, value core.Value) core.Value {
	atom.Ref = value
	return value
}

// swap updates the atom's value with the result of applying the function to the atom's current
// value followed by the remaining arguments, and returns the new value.
func swap(args ...core.Value) core.Value {
	if err := typ.ExpectAtLeastNArgs(2, args); err.NotNil() {
		return err
	}

	atom, err := typ.As[*typ.Atom](args[0])
	if err.NotNil() {
		return err
	}
	fn, err := typ.As[typ.Function](args[1])
	if err.NotNil() {
		return err
	}

	params := make([]core.Value, len(args)-1)
	params[0] = atom.Ref
	copy(params[1:], args[2:])

	// fn may be a user-defined closure, which can return a TailCall for tail-call optimization
	// instead of a realized value.
	// It must be resolved (trampolined) here, otherwise an unevaluated TailCall would be stored in
	// the atom and leak to the caller.
	res := typ.EnsureResolved(fn(params...))
	if ut.Is[typ.Fatal](res) {
		return res
	}
	atom.Ref = res
	return res
}
