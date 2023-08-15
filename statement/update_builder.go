package statement

import (
	"database/sql"
	"strings"

	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/maddiesch/go-raptor/statement/dialect"
	"github.com/maddiesch/go-raptor/statement/generator"
	"github.com/maddiesch/go-raptor/statement/query"
)

type UpdateBuilder struct {
	table     string
	set       []UpdateValue
	where     conditional.Conditional
	returning []UpdateReturnValue
}

type UpdateValue struct {
	ColumnName string
	Value      any
}

type UpdateReturnValue struct {
	ColumnName  string
	ColumnAlias string
}

func Update(t string) *UpdateBuilder {
	return &UpdateBuilder{
		table: t,
	}
}

func (b *UpdateBuilder) Set(set ...UpdateValue) *UpdateBuilder {
	b.set = append(b.set, set...)
	return b
}

func (b *UpdateBuilder) SetValue(c string, v any) *UpdateBuilder {
	return b.Set(UpdateValue{c, v})
}

func (b *UpdateBuilder) SetMap(m map[string]any) *UpdateBuilder {
	v := make([]UpdateValue, 0, len(m))
	for k, val := range m {
		v = append(v, UpdateValue{
			ColumnName: k,
			Value:      val,
		})
	}
	return b.Set(v...)
}

func (b *UpdateBuilder) Where(c conditional.Conditional) *UpdateBuilder {
	b.where = c
	return b
}

func (b *UpdateBuilder) ReturningColumn(c ...string) *UpdateBuilder {
	for _, col := range c {
		b.returning = append(b.returning, UpdateReturnValue{
			ColumnName: col,
		})
	}
	return b
}

func (b *UpdateBuilder) Returning(c ...UpdateReturnValue) *UpdateBuilder {
	b.returning = append(b.returning, c...)
	return b
}

func (b *UpdateBuilder) Generate() (string, []any, error) {
	var query query.Builder
	var args []any

	query.WriteString("UPDATE " + dialect.Identifier(b.table) + " SET")

	provider := generator.NewIncrementingArgumentNameProvider()

	var values []string
	for _, up := range b.set {
		vName := provider.Next()
		values = append(values, dialect.Identifier(up.ColumnName)+" = $"+vName)
		args = append(args, sql.Named(vName, up.Value))
	}
	if len(values) > 0 {
		query.WriteString(" " + strings.Join(values, ", "))
	}

	if b.where != nil {
		q, wArgs := b.where.Generate(provider)
		query.WriteString(" WHERE " + q)
		args = append(args, wArgs...)
	}

	if len(b.returning) > 0 {
		var columns []string
		for _, c := range b.returning {
			var s string
			if c.ColumnAlias != "" {
				s = c.ColumnName + " AS " + dialect.Identifier(c.ColumnAlias)
			} else {
				s = dialect.Identifier(c.ColumnName)
			}
			columns = append(columns, s)
		}
		query.WriteString(" RETURNING " + strings.Join(columns, ", "))
	}

	return query.String(), args, nil
}

var _ generator.Generator = (*UpdateBuilder)(nil)
