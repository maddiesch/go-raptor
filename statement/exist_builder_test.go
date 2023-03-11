package statement_test

import (
	"database/sql"
	"testing"

	"github.com/maddiesch/go-raptor/statement"
	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExists(t *testing.T) {
	sel := statement.Select().From("Foo").Where(conditional.Equal("Bar", 1))

	query, args, err := statement.Exists(sel).Generate()
	require.NoError(t, err)

	assert.Equal(t, `SELECT EXISTS(SELECT * FROM "Foo" WHERE "Bar" = $v1);`, query)

	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", 1), args[0])
	}
}
