package envs

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/fun"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

// MalFun holds the built-in Mal functions.
var MalFun core.Env = core.Bindings{
	// Atoms.
	"atom":   typ.I(atom),
	"atom?":  typ.I(isX[*typ.Atom]),
	"deref":  typ.I(deref),
	"reset!": typ.II(reset),
	"swap!":  typ.Function(swap),

	// Containers.
	"concat": typ.Function(concat),
	"cons":   typ.II(cons),
	"count":  typ.I(count),
	"empty?": typ.I(isEmpty),
	"first":  typ.I(first),
	"list":   typ.Function(list),
	"list?":  typ.I(isX[typ.List]),
	"nth":    typ.II(nth),
	"rest":   typ.I(rest),
	"vec":    typ.I(vec),

	// Fun.
	"apply": typ.Function(apply),
	"map":   typ.II(mapfn),

	// Strings.
	"pr-str":  typ.Function(prStr),
	"println": typ.Function(prnln),
	"prn":     typ.Function(prn),
	"str":     typ.Function(str),

	// Other.
	"=":           typ.II(equal),
	"macro?":      typ.I(isX[typ.Macro]),
	"read-string": typ.I(readString),
	"slurp":       typ.I(slurp),
	"throw":       typ.I(typ.Exception),
}.AsEnv()

///////////
// Atoms //

// atom creates a new atom holding the given value.
func atom(arg core.Value) *typ.Atom { return &typ.Atom{Ref: arg} }

// deref returns the value currently held by the atom.
func deref(atom *typ.Atom) core.Value { return atom.Ref }

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

////////////////
// Containers //

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

// equal returns true when two values are equal.
func equal(left, right core.Value) typ.Bool { return typ.Bool(left.Equal(right)) }

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

// isEmpty returns true when the given container has no elements.
func isEmpty(arg core.Value) core.Value {
	size, err := length(arg)
	if err.NotNil() {
		return err
	}
	return typ.Bool(size == 0)
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

// list creates a new list from the given arguments.
func list(args ...core.Value) core.Value { return typ.List(args) }

// nth returns the nth element of the given list/vector.
func nth(container core.Value, index typ.Integer) core.Value {
	slice, err := typ.AsListOrVec(container)
	if err.NotNil() {
		return err
	}

	if 0 > index || int(index) >= len(slice) {
		return typ.Exceptionf("index %d out of bounds for slice of size %d", index, len(slice))
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

// vec converts a list to a vector, returns vectors as is and errors on other types.
func vec(arg core.Value) core.Value {
	switch typed := arg.(type) {
	case typ.List:
		return typ.Vector(typed)
	case typ.Vector:
		return typed
	default:
		return typ.Expected(core.VECTOR+" or "+core.LIST, typed)
	}
}

/////////
// Fun //

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

/////////////////////
// Print functions //

// prnln prints the raw representations of the arguments separated by spaces and followed by a
// newline.
func prnln(args ...core.Value) core.Value {
	fmt.Println(strings.Join(fun.Map(rawString, args), " "))
	return typ.Nil{}
}

// prn prints the string representations of the arguments separated by spaces, followed by a
// newline.
func prn(args ...core.Value) core.Value {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}
	fmt.Print("\n")

	return typ.Nil{}
}

// prStr returns the string representations of the arguments joined by spaces.
func prStr(args ...core.Value) core.Value {
	return typ.String(strings.Join(fun.Map(core.Value.String, args), " "))
}

// rawString returns the raw (unescaped) string representation of a value.
func rawString(val core.Value) string {
	switch value := val.(type) {
	case typ.String:
		return string(value)
	case typ.List:
		return "(" + strings.Join(fun.Map(rawString, value), " ") + ")"
	case typ.Vector:
		return "[" + strings.Join(fun.Map(rawString, value), " ") + "]"
	}
	return val.String()
}

// str concatenates the raw representations of the arguments into a single string.
func str(args ...core.Value) core.Value {
	return typ.String(strings.Join(fun.Map(rawString, args), ""))
}

///////////
// Other //

// isX returns true when the type of arg is X.
func isX[X core.Value](arg core.Value) typ.Bool { return typ.Bool(ut.Is[X](arg)) }

// readString parses the first form of the given source string and returns it.
func readString(source typ.String) core.Value {
	res, err := reader.ParseOne(string(source))
	if err != nil {
		return typ.Exception(typ.String(err.Error()))
	}

	return res
}

// slurp reads the whole file at the given path and returns its contents as a string.
func slurp(path typ.String) core.Value {
	res, err := os.ReadFile(string(path))
	if err != nil {
		return typ.Exception(typ.String(err.Error()))
	}

	return typ.String(res)
}
