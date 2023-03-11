package dialect

import (
	"strconv"
)

func Identifier(n string) string {
	return strconv.Quote(n)
}
