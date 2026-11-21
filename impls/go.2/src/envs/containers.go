package envs

import (
	"slices"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/damn"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

// concat concatenates 0 or more lists.
func concat(lists ...core.Value) core.Value {
	typed := make([][]core.Value, len(lists))
	for i, elt := range lists {
		ellist, err := typ.AsListOrVec(elt)
		if err.NotNil() {
			return err
		}
		typed[i] = ellist
	}

	return typ.List(slices.Concat(typed...))
}

// cons creates a new list containing value and all the elements of list.
func cons(value core.Value, list core.Value) core.Value {
	typlist, err := typ.AsListOrVec(list)
	if err.NotNil() {
		return err
	}

	res := make(typ.List, len(typlist)+1)
	res[0] = value
	copy(res[1:], typlist)
	return res
}

// count returns the number of elements in the given container.
func count(arg core.Value) core.Value {
	res, err := length(arg)
	if err.NotNil() {
		return err
	}
	return typ.Integer(res)
}

// first returns the first element of a list/vector, nil when empty.
func first(container core.Value) core.Value {
	if ut.Is[typ.Nil](container) {
		return container
	}

	slice, err := typ.AsListOrVec(container)
	if err.NotNil() {
		return err
	}

	if len(slice) == 0 {
		return typ.Nil{}
	}

	return slice[0]
}

// length returns the number of elements in the given container, or an error if the value is not a
// container.
func length(container core.Value) (res int, err typ.Fatal) {
	switch typed := container.(type) {
	case typ.Nil:
		res = 0
	case typ.List:
		res = len(typed)
	case typ.Map:
		res = len(typed)
	case typ.MapForm:
		res = len(typed)
	case typ.Vector:
		res = len(typed)
	case typ.VectorForm:
		res = len(typed)
	default:
		err = typ.Expected("a container", container)
	}
	return
}

// list2vec converts a list to a vector, returns vectors as is and errors on other types.
func list2vec(arg core.Value) core.Value {
	switch typed := arg.(type) {
	case typ.List:
		return typ.Vector(typed)
	case typ.Vector:
		return typed
	default:
		return typ.Expected(core.VECTOR+" or "+core.LIST, typed)
	}
}

// mkList creates a new list from a list of arguments.
func mkList(args ...core.Value) core.Value { return typ.List(args) }

// mkVector creates a new vector from a list of arguments.
func mkVector(args ...core.Value) core.Value { return typ.Vector(args) }

// nth returns the nth element of the given list/vector.
func nth(container core.Value, index typ.Integer) core.Value {
	slice, err := typ.AsListOrVec(container)
	if err.NotNil() {
		return err
	}

	if 0 > index || int(index) >= len(slice) {
		return typ.ThrowErr(&damn.BoundsError{Index: int(index), Len: len(slice)})
	}
	return slice[index]
}

// rest returns the given list/vector without the first element.
func rest(container core.Value) core.Value {
	if ut.Is[typ.Nil](container) {
		return typ.List{}
	}

	slice, err := typ.AsListOrVec(container)
	if err.NotNil() {
		return err
	}

	if len(slice) == 0 {
		return typ.List{}
	}
	return typ.List(slice[1:])
}
