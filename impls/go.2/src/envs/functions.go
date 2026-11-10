package envs

import (
	"fmt"
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/fun"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

var MalFunctions core.Env = core.Bindings{
	"list":   typ.Function(list),
	"list?":  typ.I(isList),
	"empty?": typ.I(isEmpty),
	"count":  typ.I(count),
	"=":      typ.II(equal),

	"pr-str":  typ.Function(prStr),
	"println": typ.Function(prnln),
	"prn":     typ.Function(prn),
	"str":     typ.Function(str),
}.AsEnv()

////////////////
// Containers //

func count(arg core.Value) core.Value {
	res, err := length(arg)
	if err.NotNil() {
		return err
	}
	return typ.Integer(res)
}

func equal(left, right core.Value) typ.Bool {
	return typ.Bool(left.Equal(right))
}

func isEmpty(arg core.Value) core.Value {
	size, err := length(arg)
	if err.NotNil() {
		return err
	}
	return typ.Bool(size == 0)
}

func isList(arg core.Value) typ.Bool {
	return typ.Bool(ut.Is[typ.List](arg))
}

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

func list(args ...core.Value) core.Value { return typ.List(args) }

/////////////////////
// Print functions //

func prnln(args ...core.Value) core.Value {
	fmt.Println(strings.Join(fun.Map(rawString, args), " "))
	return typ.Nil{}
}

func prn(args ...core.Value) core.Value {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}

	return typ.Nil{}
}

func prStr(args ...core.Value) core.Value {
	return typ.String(strings.Join(fun.Map(core.Value.String, args), " "))
}

func rawString(val core.Value) string {
	if str, isStr := val.(typ.String); isStr {
		return string(str)
	}
	return val.String()
}

func str(args ...core.Value) core.Value {
	return typ.String(strings.Join(fun.Map(rawString, args), ""))
}
