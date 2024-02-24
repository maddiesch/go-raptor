package statement

import (
	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/maddiesch/go-raptor/statement/dialect"
	"github.com/maddiesch/go-raptor/statement/generator"
	"github.com/maddiesch/go-raptor/statement/query"
)

func Delete() *DeleteBuilder {
	return &DeleteBuilder{}
}

type DeleteBuilder struct {
	tableName string
	where     conditional.Conditional
}

func (b *DeleteBuilder) From(tableName string) *DeleteBuilder {
	b.tableName = tableName

	return b
}

func (b *DeleteBuilder) Where(where conditional.Conditional) *DeleteBuilder {
	b.where = where

	return b
}

func (b *DeleteBuilder) Generate() (string, []any, error) {
	var query query.Builder
	var args []any

	_, _ = query.WriteStringf("DELETE FROM %s", dialect.Identifier(b.tableName))

	provider := generator.NewIncrementingArgumentNameProvider()

	if b.where != nil {
		where, wArgs := b.where.Generate(provider)
		args = append(args, wArgs...)
		_, _ = query.WriteStringf(" WHERE %s", where)
	}

	return query.String(), args, nil
}
