package envs

import (
	"os"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/typ"
)

// equal returns true when two values are equal.
func equal(left, right core.Value) typ.Bool { return typ.Bool(left.Equal(right)) }

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
	if str, ok := arg.(interface{ Unbox() string }); ok {
		return str.Unbox(), typ.Fatal{}
	}
	return "", typ.Expected("keyword, string or symbol", arg)
}

// mkKeyword creates a new keyword from a keyword, string or symbol.
func mkKeyword(name core.Value) core.Value {
	namestr, err := goString(name)
	if err.NotNil() {
		return err
	}
	return typ.NewKeyword(namestr)
}

// mkSymbol creates a new symbol from a keyword, string or symbol.
func mkSymbol(name core.Value) core.Value {
	namestr, err := goString(name)
	if err.NotNil() {
		return err
	}
	return typ.NewSymbol(namestr)
}

// readString parses the first form of the given source string and returns it.
func readString(source typ.String) core.Value {
	res, err := reader.ParseOne(source.Unbox())
	if err != nil {
		return typ.Throw(typ.NewString(err.Error()))
	}

	return res
}

// slurp reads the whole file at the given path and returns its contents as a string.
func slurp(path typ.String) core.Value {
	res, err := os.ReadFile(path.Unbox())
	if err != nil {
		return typ.Throw(typ.NewString(err.Error()))
	}

	return typ.NewString(string(res))
}
