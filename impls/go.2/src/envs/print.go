package envs

import (
	"fmt"
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/fun"
	"github.com/mooss/mal/src/typ"
)

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
	return typ.NewString(strings.Join(fun.Map(core.Value.String, args), " "))
}

// rawString returns the raw (unescaped) string representation of a value.
func rawString(val core.Value) string {
	switch value := val.(type) {
	case typ.String:
		return value.Unbox()
	case typ.List:
		return "(" + strings.Join(fun.Map(rawString, value), " ") + ")"
	case typ.Vector:
		return "[" + strings.Join(fun.Map(rawString, value), " ") + "]"
	}
	return val.String()
}

// str concatenates the raw representations of the arguments into a single string.
func str(args ...core.Value) core.Value {
	return typ.NewString(strings.Join(fun.Map(rawString, args), ""))
}
