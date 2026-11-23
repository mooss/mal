package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

var OddMapArgsErr = typ.Throwf("odd-numbered map arguments")

// MalFun holds the built-in Mal functions.
var MalFun core.Env = core.Bindings{
	// Atoms.
	"atom":   typ.I(mkAtom),
	"deref":  typ.I(deref),
	"reset!": typ.II(reset),
	"swap!":  typ.Function(swap),

	// Containers.
	"concat": typ.Function(concat),
	"cons":   typ.II(cons),
	"count":  typ.I(count),
	"first":  typ.I(first),
	"list":   typ.Function(mkList),
	"nth":    typ.II(nth),
	"rest":   typ.I(rest),
	"vec":    typ.I(list2vec),
	"vector": typ.Function(mkVector),

	// Higher-order.
	"apply": typ.Function(apply),
	"map":   typ.II(mapfn),

	// Maps.
	"assoc":     typ.Function(assoc),
	"contains?": typ.II(typ.Map.MalContains),
	"dissoc":    typ.Function(dissoc),
	"get":       typ.II(get),
	"hash-map":  typ.Function(mkMap),
	"keys":      typ.I(typ.Map.MalKeys),
	"vals":      typ.I(typ.Map.MalValues),

	// Predicates.
	"atom?":       typ.I(isX[*typ.Atom]),
	"bool?":       typ.I(isX[typ.Bool]),
	"empty?":      typ.I(isEmpty),
	"false?":      typ.I(isFalse),
	"fn?":         typ.I(isX[typ.Function]),
	"map?":        typ.I(isX[typ.Map]),
	"keyword?":    typ.I(isX[typ.Keyword]),
	"list?":       typ.I(isX[typ.List]),
	"macro?":      typ.I(isX[typ.Macro]),
	"nil?":        typ.I(isX[typ.Nil]),
	"number?":     typ.I(isX[typ.Integer]),
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
	"keyword":     typ.I(mkKeyword),
	"read-string": typ.I(readString),
	"readline":    typ.I(readline),
	"slurp":       typ.I(slurp),
	"symbol":      typ.I(mkSymbol),
	"throw":       typ.I(typ.Throw),

	// NIY.
	"time-ms":   typ.Function(nimpl),
	"meta":      typ.Function(nimpl),
	"with-meta": typ.Function(nimpl),
	"seq":       typ.Function(nimpl),
	"conj":      typ.Function(nimpl),
}.AsEnv()

func nimpl(...core.Value) core.Value { return typ.Throwf("not implemented") }
