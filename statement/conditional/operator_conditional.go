package conditional

import (
	"database/sql"
	"fmt"

	"github.com/maddiesch/go-raptor/statement/dialect"
	"github.com/maddiesch/go-raptor/statement/generator"
)

func Equal(col string, val any) Conditional {
	if val == nil {
		return Null(col)
	}
	return &operatorInfixConditional{col, "=", val}
}

func NotEqual(col string, val any) Conditional {
	if val == nil {
		return NotNull(col)
	}
	return &operatorInfixConditional{col, "!=", val}
}

func LessThan(col string, val any) Conditional {
	return &operatorInfixConditional{col, "<", val}
}

func LessThanEq(col string, val any) Conditional {
	return &operatorInfixConditional{col, "<=", val}
}

func GreaterThan(col string, val any) Conditional {
	return &operatorInfixConditional{col, ">", val}
}

func GreaterThanEq(col string, val any) Conditional {
	return &operatorInfixConditional{col, ">=", val}
}

type operatorInfixConditional struct {
	column   string
	operator string
	value    any
}

func (c *operatorInfixConditional) Generate(provider generator.ArgumentNameProvider) (string, []any) {
	name := provider.Next()

	return fmt.Sprintf("%s %s $%s", dialect.Identifier(c.column), c.operator, name), []any{sql.Named(name, c.value)}
}

func Null(col string) Conditional {
	return &nullConditional{col, true}
}

func NotNull(col string) Conditional {
	return &nullConditional{col, false}
}

type nullConditional struct {
	column string
	isNull bool
}

func (c *nullConditional) Generate(provider generator.ArgumentNameProvider) (string, []any) {
	v := "NOT NULL"
	if c.isNull {
		v = "NULL"
	}
	return fmt.Sprintf("%s IS %s", dialect.Identifier(c.column), v), nil
}
