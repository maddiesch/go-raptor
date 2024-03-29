package statement

import (
	"strings"

	"github.com/maddiesch/go-raptor/statement/dialect"
	"github.com/maddiesch/go-raptor/statement/query"
)

func CreateTable(name string) *CreateTableBuilder {
	return &CreateTableBuilder{tableName: name}
}

func Column(name string, cType ColumnType) *ColumnBuilder {
	return &ColumnBuilder{name: name, cType: cType, nullable: true}
}

type ColumnType string

const (
	ColumnTypeText    ColumnType = "TEXT"
	ColumnTypeInteger ColumnType = "INTEGER"
	ColumnTypeFloat   ColumnType = "REAL"
	ColumnTypeBlob    ColumnType = "BLOB"
)

type ColumnBuilder struct {
	name           string
	defaultLiteral string
	cType          ColumnType
	nullable       bool
	unique         bool
	pk             bool
}

func (c *ColumnBuilder) NotNull() *ColumnBuilder {
	c.nullable = false

	return c
}

func (c *ColumnBuilder) Unique() *ColumnBuilder {
	c.unique = true

	return c
}

func (c *ColumnBuilder) Default(literal string) *ColumnBuilder {
	c.defaultLiteral = literal

	return c
}

func (c *ColumnBuilder) Generate() (string, []any, error) {
	var q query.Builder
	_, _ = q.WriteString(dialect.Identifier(c.name))
	_, _ = q.WriteRune(' ')
	_, _ = q.WriteString(string(c.cType))

	if c.pk {
		_, _ = q.WriteString(" PRIMARY KEY")
	}

	if !c.nullable {
		_, _ = q.WriteString(" NOT NULL")
	}
	if c.unique {
		_, _ = q.WriteString(" UNIQUE")
	}
	if c.defaultLiteral != "" {
		_, _ = q.WriteString(" DEFAULT ")
		_, _ = q.WriteString(c.defaultLiteral)
	}

	return q.Builder.String(), nil, nil
}

type CreateTableBuilder struct {
	tableName   string
	primaryKey  *ColumnBuilder
	columns     []*ColumnBuilder
	ifNotExists bool
}

func (c *CreateTableBuilder) IfNotExists() *CreateTableBuilder {
	c.ifNotExists = true

	return c
}

func (c *CreateTableBuilder) PrimaryKey(name string, cType ColumnType) *CreateTableBuilder {
	c.primaryKey = &ColumnBuilder{name: name, cType: cType, nullable: false, unique: true, pk: true}

	return c
}

func (c *CreateTableBuilder) Column(column ...*ColumnBuilder) *CreateTableBuilder {
	c.columns = append(c.columns, column...)

	return c
}

func (c *CreateTableBuilder) Generate() (string, []any, error) {
	var query query.Builder

	_, _ = query.WriteString("CREATE TABLE")
	if c.ifNotExists {
		_, _ = query.WriteString(" IF NOT EXISTS")
	}
	_, _ = query.WriteStringf(" %s", dialect.Identifier(c.tableName))
	_, _ = query.WriteString(" (")

	var args []any
	var columns []string

	if c.primaryKey != nil {
		sub, sArgs, err := c.primaryKey.Generate()
		if err != nil {
			return "", nil, err
		}
		columns = append(columns, sub)
		args = append(args, sArgs...)
	}

	for _, column := range c.columns {
		sub, sArgs, err := column.Generate()
		if err != nil {
			return "", nil, err
		}
		columns = append(columns, sub)
		args = append(args, sArgs...)
	}

	_, _ = query.WriteString(strings.Join(columns, ", "))

	_, _ = query.WriteString(")")

	return query.String(), args, nil
}
