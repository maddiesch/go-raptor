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
	orIgnore  bool
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
	b.orIgnore = false

	return b
}

func (b *InsertBuilder) OrIgnore() *InsertBuilder {
	b.orIgnore = true
	b.orReplace = false

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
	if b.orIgnore {
		query.WriteString("OR IGNORE ")
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

func (b *InsertBuilder) Returning(col ...string) *InsertReturnBuilder {
	return &InsertReturnBuilder{
		Insert:  b,
		Columns: col,
	}
}

type InsertReturnBuilder struct {
	Insert  *InsertBuilder
	Columns []string
}

func (b *InsertReturnBuilder) Generate() (string, []any, error) {
	q, args, err := b.Insert.Generate()
	if err != nil {
		return "", nil, err
	}

	var columns []string
	for _, c := range b.Columns {
		columns = append(columns, dialect.Identifier(c))
	}
	if len(columns) == 0 {
		columns = append(columns, "*")
	}

	return q[:len(q)-1] + " RETURNING " + strings.Join(columns, ", ") + ";", args, nil
}
