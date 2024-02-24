package statement_test

import (
	"database/sql"
	"testing"

	"github.com/maddiesch/go-raptor/statement"
	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/maddiesch/go-raptor/statement/generator"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBuilder(t *testing.T) {
	tests := []struct {
		statement     generator.Generator
		expectedQuery string
		expectedArgs  []any
	}{
		{
			statement:     statement.Update("testing").SetValue("name", "Maddie"),
			expectedQuery: `UPDATE "testing" SET "name" = $v1;`,
			expectedArgs:  []any{sql.Named("v1", "Maddie")},
		},
		{
			statement:     statement.Update("testing").SetValue("name", "Maddie").Where(conditional.Equal("id", 1)),
			expectedQuery: `UPDATE "testing" SET "name" = $v1 WHERE "id" = $v2;`,
			expectedArgs:  []any{sql.Named("v1", "Maddie"), sql.Named("v2", 1)},
		},
		{
			statement:     statement.Update("testing").SetValue("name", "Maddie").Where(conditional.Equal("id", 1)).ReturningColumn("updated_at"),
			expectedQuery: `UPDATE "testing" SET "name" = $v1 WHERE "id" = $v2 RETURNING "updated_at";`,
			expectedArgs:  []any{sql.Named("v1", "Maddie"), sql.Named("v2", 1)},
		},
		{
			statement:     statement.Update("testing").SetValue("name", "Maddie").Where(conditional.Equal("id", 1)).ReturningColumn("updated_at").Returning(statement.UpdateReturnValue{ColumnName: "1", ColumnAlias: "ReturnValue"}),
			expectedQuery: `UPDATE "testing" SET "name" = $v1 WHERE "id" = $v2 RETURNING "updated_at", 1 AS "ReturnValue";`,
			expectedArgs:  []any{sql.Named("v1", "Maddie"), sql.Named("v2", 1)},
		},
		{
			statement:     statement.Update("testing").SetMap(map[string]any{"name": "Maddie"}).Where(conditional.Equal("id", 1)).ReturningColumn("updated_at").Returning(statement.UpdateReturnValue{ColumnName: "1", ColumnAlias: "ReturnValue"}),
			expectedQuery: `UPDATE "testing" SET "name" = $v1 WHERE "id" = $v2 RETURNING "updated_at", 1 AS "ReturnValue";`,
			expectedArgs:  []any{sql.Named("v1", "Maddie"), sql.Named("v2", 1)},
		},
	}

	for _, test := range tests {
		t.Run(test.expectedQuery, func(t *testing.T) {
			query, args, err := test.statement.Generate()
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectedQuery, query)
				assert.Equal(t, test.expectedArgs, args)
			}
		})
	}
}
