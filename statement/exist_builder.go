package statement

import (
	"fmt"

	"github.com/maddiesch/go-raptor/statement/generator"
)

func Exists(sel *SelectBuilder) generator.Generator {
	return &existsQueryBuilder{sel}
}

type existsQueryBuilder struct {
	child *SelectBuilder
}

func (b *existsQueryBuilder) Generate() (string, []any, error) {
	child, args, err := b.child.Generate()

	return fmt.Sprintf("SELECT EXISTS(%s);", child[:len(child)-1]), args, err
}
