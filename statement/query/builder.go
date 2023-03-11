package query

import (
	"fmt"
	"strings"
)

type Builder struct {
	strings.Builder
}

func (b *Builder) String() string {
	return b.Builder.String() + ";"
}

func (b *Builder) WriteStringf(format string, a ...any) (int, error) {
	return b.Builder.WriteString(fmt.Sprintf(format, a...))
}
