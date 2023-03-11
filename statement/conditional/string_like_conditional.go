package conditional

import (
	"database/sql"
	"fmt"

	"github.com/maddiesch/go-raptor/statement/dialect"
	"github.com/maddiesch/go-raptor/statement/generator"
)

func StringLike(column, value string) Conditional {
	return &stringLikeConditional{column, value}
}

type stringLikeConditional struct {
	column string
	value  string
}

func (c *stringLikeConditional) Generate(p generator.ArgumentNameProvider) (string, []any) {
	name := p.Next()

	return fmt.Sprintf("%s LIKE $%s", dialect.Identifier(c.column), name), []any{sql.Named(name, c.value)}
}

func StringHasPrefix(column, value string) Conditional {
	return StringLike(column, value+`%`)
}

func StringHasSuffix(column, value string) Conditional {
	return StringLike(column, `%`+value)
}

func CaseInsensitive(c Conditional) Conditional {
	return &caseInsensitiveConditional{c}
}

type caseInsensitiveConditional struct {
	child Conditional
}

func (c *caseInsensitiveConditional) Generate(p generator.ArgumentNameProvider) (string, []any) {
	child, args := c.child.Generate(p)

	return child + " COLLATE NOCASE", args
}
