package envs

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/damn"
	"github.com/mooss/mal/src/fun"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

var OddMapArgsErr = typ.Throwf("odd-numbered map arguments")

// MalFun holds the built-in Mal functions.
var MalFun core.Env = core.Bindings{
	// Atoms.
	"deref":  typ.I(deref),
	"reset!": typ.II(reset),
	"swap!":  typ.Function(swap),

	// Constructors.
	"atom":     typ.I(mkAtom),
	"hash-map": typ.Function(mkMap),
	"symbol":   typ.I(mkSymbol),
	"keyword":  typ.I(mkKeyword),
	"vector":   typ.Function(mkVector),

	// Containers.
	"concat": typ.Function(concat),
	"cons":   typ.II(cons),
	"count":  typ.I(count),
	"first":  typ.I(first),
	"list":   typ.Function(mkList),
	"nth":    typ.II(nth),
	"rest":   typ.I(rest),
	"vec":    typ.I(vec),

	// Fun.
	"apply": typ.Function(apply),
	"map":   typ.II(mapfn),

	// Maps.
	"assoc":     typ.Function(assoc),
	"contains?": typ.II(typ.Map.MalContains),
	"dissoc":    typ.Function(dissoc),
	"get":       typ.II(get),
	"keys":      typ.I(typ.Map.MalKeys),
	"vals":      typ.I(typ.Map.MalValues),

	// Predicates.
	"atom?":       typ.I(isX[*typ.Atom]),
	"bool?":       typ.I(isX[typ.Bool]),
	"empty?":      typ.I(isEmpty),
	"false?":      typ.I(isFalse),
	"map?":        typ.I(isX[typ.Map]),
	"int?":        typ.I(isX[typ.Integer]),
	"keyword?":    typ.I(isX[typ.Keyword]),
	"list?":       typ.I(isX[typ.List]),
	"macro?":      typ.I(isX[typ.Macro]),
	"nil?":        typ.I(isX[typ.Nil]),
	"sequential?": typ.I(isSequential),
	"string?":     typ.I(isX[typ.String]),
	"symbol?":     typ.I(isX[typ.Symbol]),
	"true?":       typ.I(isTrue),
	"vector?":     typ.I(isX[typ.Vector]),

	// Strings.
	"pr-str":  typ.Function(prStr),
	"println": typ.Function(prnln),
	"prn":     typ.Function(prn),
	"str":     typ.Function(str),

	// Other.
	"=":           typ.II(equal),
	"read-string": typ.I(readString),
	"slurp":       typ.I(slurp),
	"throw":       typ.I(typ.Throw),
}.AsEnv()

///////////
// Atoms //

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

//////////////////
// Constructors //

// mkAtom creates a new atom from a value.
func mkAtom(arg core.Value) *typ.Atom { return &typ.Atom{Ref: arg} }

// mkKeyword creates a new keyword from a keyword, string or symbol.
func mkKeyword(name core.Value) core.Value {
	namestr, err := goString(name)
	if err.NotNil() {
		return err
	}
	return typ.Keyword(namestr)
}

// mkList creates a new list from a list of arguments.
func mkList(args ...core.Value) core.Value { return typ.List(args) }

// mkMap creates a new map from even key, value... arguments.
func mkMap(args ...core.Value) core.Value {
	if len(args)%2 != 0 {
		return OddMapArgsErr
	}

	return extendMap(make(typ.Map, len(args)/2), args...)
}

// mkSymbol creates a new symbol from a keyword, string or symbol.
func mkSymbol(name core.Value) core.Value {
	namestr, err := goString(name)
	if err.NotNil() {
		return err
	}
	return typ.Symbol(namestr)
}

// mkVector creates a new vector from a list of arguments.
func mkVector(args ...core.Value) core.Value { return typ.Vector(args) }

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

//////////
// Maps //

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

////////////////
// Predicates //

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

// extendMap adds new associations to dest from even key, value... arguments.
// Will crash when len(args) is not even.
func extendMap(dest typ.Map, args ...core.Value) core.Value {
	for i := 0; i < len(args); i += 2 {
		if !typ.ValidMapKey(args[i]) {
			return typ.Throwf("invalid map key %#v (%T)", args[i], args[i])
		}
		dest[args[i]] = args[i+1]
	}
	return dest
}

// goString returns the raw Go string underlying a Keyword, String or Symbol.
func goString(arg core.Value) (string, typ.Fatal) {
	switch typed := arg.(type) {
	case typ.Keyword:
		return string(typed), typ.Fatal{}
	case typ.String:
		return string(typed), typ.Fatal{}
	case typ.Symbol:
		return string(typed), typ.Fatal{}
	default:
		return "", typ.Expected("keyword, string or symbol", arg)
	}
}

// readString parses the first form of the given source string and returns it.
func readString(source typ.String) core.Value {
	res, err := reader.ParseOne(string(source))
	if err != nil {
		return typ.Throw(typ.String(err.Error()))
	}

	return res
}

// slurp reads the whole file at the given path and returns its contents as a string.
func slurp(path typ.String) core.Value {
	res, err := os.ReadFile(string(path))
	if err != nil {
		return typ.Throw(typ.String(err.Error()))
	}

	return typ.String(res)
}
