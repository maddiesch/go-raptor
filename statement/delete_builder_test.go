package statement_test

import (
	"database/sql"
	"testing"

	"github.com/maddiesch/go-raptor/statement"
	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteBuilder(t *testing.T) {
	t.Run("basic query generator", func(t *testing.T) {
		query, args, err := statement.Delete().From("TestTable").Where(conditional.Equal("Foo", "Bar")).Generate()

		require.NoError(t, err)

		assert.Equal(t, `DELETE FROM "TestTable" WHERE "Foo" = $v1;`, query)

		if assert.Len(t, args, 1) {
			assert.Equal(t, sql.Named("v1", "Bar"), args[0])
		}
	})
}
