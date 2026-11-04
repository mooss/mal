package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mooss/mal/src/core"
)

func PrintStr(typ core.Type) string {
	switch value := typ.(type) {
	case core.Integer:
		return strconv.FormatInt(int64(value), 10)
	case core.Symbol:
		return string(value)
	case core.List:
		return "(" + strings.Join(value.Map(PrintStr), " ") + ")"
	case core.Error:
		return fmt.Sprintf("(error %q)", value.Error())
	}
	return fmt.Sprintf("!! PrintStr: no support for %T: %#v", typ, typ)
}
