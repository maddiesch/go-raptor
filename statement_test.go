package raptor_test

import (
	"testing"

	"github.com/maddiesch/go-raptor/internal/test"
	"github.com/maddiesch/go-raptor/statement"
	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConn_QueryRowStatement(t *testing.T) {
	conn, ctx := test.Setup(t)
	defer conn.Close()

	query := statement.Select("FirstName", "LastName").From("People").Where(conditional.Equal("FirstName", "Maddie")).Limit(1)

	var firstName, lastName string
	err := conn.QueryRowStatement(ctx, query).Scan(&firstName, &lastName)

	assert.NoError(t, err)

	assert.Equal(t, "Maddie", firstName)
	assert.Equal(t, "Schipper", lastName)
}

func TestConn_QueryStatement(t *testing.T) {
	conn, ctx := test.Setup(t)
	defer conn.Close()

	query := statement.Select("FirstName", "LastName").From("People").Where(conditional.Equal("FirstName", "Maddie")).Limit(1)

	rows, err := conn.QueryStatement(ctx, query)
	require.NoError(t, err)

	assert.True(t, rows.Next())

	var firstName, lastName string
	err = rows.Scan(&firstName, &lastName)
	assert.NoError(t, err)
	assert.Equal(t, "Maddie", firstName)
	assert.Equal(t, "Schipper", lastName)

	assert.False(t, rows.Next())
}
