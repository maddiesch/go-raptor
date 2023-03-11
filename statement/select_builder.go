package statement

import (
	"fmt"
	"strings"

	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/maddiesch/go-raptor/statement/dialect"
	"github.com/maddiesch/go-raptor/statement/generator"
	"github.com/maddiesch/go-raptor/statement/query"
	"github.com/samber/lo"
)

type SelectBuilder struct {
	tableName  string
	isDistinct bool
	columns    []string
	where      conditional.Conditional
	limit      *int64
	offset     *int64
	orderBy    []OrderBy
}

type OrderBy struct {
	Column    string
	Ascending bool
}

func (o OrderBy) String() string {
	key := "DESC"
	if o.Ascending {
		key = "ASC"
	}
	return fmt.Sprintf("%s %s", dialect.Identifier(o.Column), key)
}

func Select(columns ...string) *SelectBuilder {
	return &SelectBuilder{
		columns: columns,
	}
}

func (b *SelectBuilder) From(table string) *SelectBuilder {
	b.tableName = table

	return b
}

func (b *SelectBuilder) Distinct() *SelectBuilder {
	b.isDistinct = true

	return b
}

func (b *SelectBuilder) Where(condition conditional.Conditional) *SelectBuilder {
	b.where = condition

	return b
}

func (b *SelectBuilder) Limit(l int64) *SelectBuilder {
	b.limit = lo.ToPtr(l)

	return b
}

func (b *SelectBuilder) Offset(o int64) *SelectBuilder {
	b.offset = lo.ToPtr(o)

	return b
}

func (b *SelectBuilder) OrderBy(col string, asc bool) *SelectBuilder {
	b.orderBy = append(b.orderBy, OrderBy{
		Column:    col,
		Ascending: asc,
	})

	return b
}

func (b *SelectBuilder) Generate() (string, []any, error) {
	var query query.Builder
	var args []any

	query.WriteString("SELECT ")

	if b.isDistinct {
		query.WriteString("DISTINCT ")
	}

	if len(b.columns) == 0 {
		query.WriteRune('*')
	} else {
		col := lo.Map(b.columns, func(c string, _ int) string {
			return dialect.Identifier(c)
		})
		query.WriteString(strings.Join(col, ", "))
	}

	query.WriteStringf(" FROM %s", dialect.Identifier(b.tableName))

	provider := generator.NewIncrementingArgumentNameProvider()

	if b.where != nil {
		where, wArgs := b.where.Generate(provider)

		query.WriteString(" WHERE ")
		query.WriteString(where)

		args = append(args, wArgs...)
	}

	if len(b.orderBy) > 0 {
		order := strings.Join(lo.Map(b.orderBy, func(o OrderBy, _ int) string {
			return o.String()
		}), ", ")
		query.WriteStringf(" ORDER BY %s", order)
	}

	if b.limit != nil {
		query.WriteStringf(" LIMIT %d", lo.FromPtr(b.limit))
	}

	if b.offset != nil {
		query.WriteStringf(" OFFSET %d", *b.offset)
	}

	return query.String(), args, nil
}

var _ generator.Generator = (*SelectBuilder)(nil)
