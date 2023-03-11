package statement_test

import (
	"database/sql"
	"testing"

	"github.com/maddiesch/go-raptor/statement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertBuilder(t *testing.T) {
	t.Run("simple insert", func(t *testing.T) {
		query, args, err := statement.Insert().Into("TestTable").ValueMap(map[string]any{
			"Name": "MTG",
			"Age":  30,
		}).Generate()

		require.NoError(t, err)

		assert.Equal(t, `INSERT INTO "TestTable" ("Age", "Name") VALUES ($v1, $v2);`, query)

		if assert.Len(t, args, 2) {
			assert.Equal(t, sql.Named("v1", 30), args[0])
			assert.Equal(t, sql.Named("v2", "MTG"), args[1])
		}
	})

	t.Run("insert or replace", func(t *testing.T) {
		query, args, err := statement.Insert().OrReplace().Into("TestTable").Value("Name", "MTG").Value("Age", 30).Generate()

		require.NoError(t, err)

		assert.Equal(t, `INSERT OR REPLACE INTO "TestTable" ("Age", "Name") VALUES ($v1, $v2);`, query)

		if assert.Len(t, args, 2) {
			assert.Equal(t, sql.Named("v1", 30), args[0])
			assert.Equal(t, sql.Named("v2", "MTG"), args[1])
		}
	})

	t.Run("no values", func(t *testing.T) {
		query, args, err := statement.Insert().Into("TestTable").Generate()

		require.NoError(t, err)

		assert.Equal(t, `INSERT INTO "TestTable" DEFAULT VALUES;`, query)

		assert.Len(t, args, 0)
	})
}
