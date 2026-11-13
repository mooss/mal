package envs

import (
	"fmt"
	"os"
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
	"atom?":  typ.I(isAtom),
	"deref":  typ.I(deref),
	"reset!": typ.II(reset),
	"swap!":  typ.Function(swap),

	// Containers.
	"list":   typ.Function(list),
	"list?":  typ.I(isList),
	"empty?": typ.I(isEmpty),
	"count":  typ.I(count),

	// Strings.
	"pr-str":  typ.Function(prStr),
	"println": typ.Function(prnln),
	"prn":     typ.Function(prn),
	"str":     typ.Function(str),

	// Other.
	"=":           typ.II(equal),
	"read-string": typ.I(readString),
	"slurp":       typ.I(slurp),
}.AsEnv()

///////////
// Atoms //

// atom creates a new atom holding the given value.
func atom(arg core.Value) *typ.Atom { return &typ.Atom{Ref: arg} }

// isAtom returns true when arg is an atom.
func isAtom(arg core.Value) typ.Bool { return typ.Bool(ut.Is[*typ.Atom](arg)) }

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

	res := fn(params...)
	if ut.Is[typ.Fatal](res) {
		return res
	}
	atom.Ref = res
	return res
}

////////////////
// Containers //

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

// isEmpty returns true when the given container has no elements.
func isEmpty(arg core.Value) core.Value {
	size, err := length(arg)
	if err.NotNil() {
		return err
	}
	return typ.Bool(size == 0)
}

// isList returns true when arg is a list.
func isList(arg core.Value) typ.Bool { return typ.Bool(ut.Is[typ.List](arg)) }

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

// readString parses the first form of the given source string and returns it.
func readString(source typ.String) core.Value {
	res, err := reader.ParseOne(string(source))
	if err != nil {
		return typ.NewFatal(err)
	}

	return res
}

// slurp reads the whole file at the given path and returns its contents as a string.
func slurp(path typ.String) core.Value {
	res, err := os.ReadFile(string(path))
	if err != nil {
		return typ.NewFatal(err)
	}

	return typ.String(res)
}
