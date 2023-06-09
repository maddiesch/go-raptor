package statement

import (
	"database/sql"
	"sort"
	"strings"

	"github.com/maddiesch/go-raptor/statement/dialect"
	"github.com/maddiesch/go-raptor/statement/generator"
	"github.com/maddiesch/go-raptor/statement/query"
	"github.com/samber/lo"
)

type InsertValue struct {
	ColumnName string
	Value      any
}

type InsertBuilder struct {
	tableName string
	orReplace bool
	values    map[string]any
}

func Insert() *InsertBuilder {
	return &InsertBuilder{
		values: make(map[string]any),
	}
}

func (b *InsertBuilder) Into(tableName string) *InsertBuilder {
	b.tableName = tableName

	return b
}

func (b *InsertBuilder) OrReplace() *InsertBuilder {
	b.orReplace = true

	return b
}

func (b *InsertBuilder) ValueMap(m map[string]any) *InsertBuilder {
	for k, v := range m {
		b.values[k] = v
	}

	return b
}

func (b *InsertBuilder) Value(column string, value any) *InsertBuilder {
	b.values[column] = value

	return b
}

func (b *InsertBuilder) Generate() (string, []any, error) {
	var query query.Builder
	var args []any

	query.WriteString("INSERT ")
	if b.orReplace {
		query.WriteString("OR REPLACE ")
	}
	query.WriteStringf("INTO %s ", dialect.Identifier(b.tableName))

	provider := generator.NewIncrementingArgumentNameProvider()

	if len(b.values) == 0 {
		query.WriteString("DEFAULT VALUES")
	} else {
		var columns, values []string

		sortedColumns := lo.Keys(b.values)
		sort.Strings(sortedColumns)

		for _, column := range sortedColumns {
			vName := provider.Next()
			columns = append(columns, dialect.Identifier(column))
			values = append(values, "$"+vName)
			args = append(args, sql.Named(vName, b.values[column]))
		}

		query.WriteStringf("(%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(values, ", "))
	}

	return query.String(), args, nil
}
