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
	}
	return fmt.Sprintf("!! PrintStr: no support for %#q (%T)", typ, typ)
}
