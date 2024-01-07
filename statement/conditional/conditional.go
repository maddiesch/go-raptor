package conditional

import (
	"database/sql"

	"github.com/maddiesch/go-raptor/statement/dialect"
	"github.com/maddiesch/go-raptor/statement/generator"
)

type Conditional interface {
	Generate(generator.ArgumentNameProvider) (string, []any)
}

type Value any

func ColumnName(column string) Conditional {
	return LiteralIdentifier(column)
}

type LiteralIdentifier string

func (c LiteralIdentifier) Generate(generator.ArgumentNameProvider) (string, []any) {
	return dialect.Identifier(string(c)), nil
}

func WrappedValue(value any) Conditional {
	return wrappedValue{value}
}

type wrappedValue struct {
	value any
}

func (v wrappedValue) Generate(p generator.ArgumentNameProvider) (string, []any) {
	n := p.Next()
	return "$" + n, []any{sql.Named(n, v.value)}
}
