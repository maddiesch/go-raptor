package statement_test

import (
	"database/sql"
	"testing"

	"github.com/maddiesch/go-raptor/statement"
	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/maddiesch/go-raptor/statement/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectBuilder(t *testing.T) {
	tests := []struct {
		statement     generator.Generator
		expectedQuery string
		expectedArgs  []any
	}{
		{
			statement:     statement.Select().From("TestTable"),
			expectedQuery: `SELECT * FROM "TestTable";`,
			expectedArgs:  nil,
		},
		{
			statement:     statement.Select().From("TestTable").Limit(10).Offset(100),
			expectedQuery: `SELECT * FROM "TestTable" LIMIT 10 OFFSET 100;`,
			expectedArgs:  nil,
		},
		{
			statement:     statement.Select().From("TestTable").Where(conditional.Equal("Key", "foo")),
			expectedQuery: `SELECT * FROM "TestTable" WHERE "Key" = $v1;`,
			expectedArgs:  []any{sql.Named("v1", "foo")},
		},
		{
			statement:     statement.Select("FirstName", "LastName").From("TestTable"),
			expectedQuery: `SELECT "FirstName", "LastName" FROM "TestTable";`,
			expectedArgs:  nil,
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

	t.Run("with a distinct limit", func(t *testing.T) {
		query, _, err := statement.Select().Distinct().From("TestTable").Limit(1).Generate()

		require.NoError(t, err)

		assert.Equal(t, `SELECT DISTINCT * FROM "TestTable" LIMIT 1;`, query)
	})

	t.Run("with a single order", func(t *testing.T) {
		query, _, err := statement.Select().Distinct().From("TestTable").OrderBy("Foo", true).Generate()

		require.NoError(t, err)

		assert.Equal(t, `SELECT DISTINCT * FROM "TestTable" ORDER BY "Foo" ASC;`, query)
	})

	t.Run("with multi-order", func(t *testing.T) {
		query, _, err := statement.Select().Distinct().From("TestTable").OrderBy("Foo", true).OrderBy("Bar", false).Generate()

		require.NoError(t, err)

		assert.Equal(t, `SELECT DISTINCT * FROM "TestTable" ORDER BY "Foo" ASC, "Bar" DESC;`, query)
	})
}
