package main

import (
	"fmt"
	"strconv"
	"strings"
)

func PrintStr(typ MalType) string {
	switch value := typ.(type) {
	case Integer:
		return strconv.FormatInt(int64(value), 10)
	case Symbol:
		return string(value)
	case List:
		return "(" + strings.Join(value.Map(PrintStr), " ") + ")"
	}
	return fmt.Sprintf("!! PrintStr: no support for %#q (%T)", typ, typ)
}
