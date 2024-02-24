package raptor_test

import (
	"context"
	"testing"

	"github.com/maddiesch/go-raptor"
	"github.com/maddiesch/go-raptor/internal/test"
	"github.com/maddiesch/go-raptor/raptortest"
	"github.com/maddiesch/go-raptor/statement"
	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConn_QueryRowStatement(t *testing.T) {
	conn, ctx := test.Setup(t)

	query := statement.Select("FirstName", "LastName").From("People").Where(conditional.Equal("FirstName", "Maddie")).Limit(1)

	var firstName, lastName string
	err := conn.QueryRowStatement(ctx, query).Scan(&firstName, &lastName)

	assert.NoError(t, err)

	assert.Equal(t, "Maddie", firstName)
	assert.Equal(t, "Schipper", lastName)
}

func TestConn_QueryStatement(t *testing.T) {
	conn, ctx := test.Setup(t)

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

func TestConn_ExecStatement(t *testing.T) {
	conn, ctx := test.Setup(t)

	exec := statement.Insert().Into("People").ValueMap(map[string]any{
		"FirstName": "Taylor",
		"LastName":  "Swift",
	})

	_, err := conn.ExecStatement(ctx, exec)
	assert.NoError(t, err)
}

func TestQueryRowStatement(t *testing.T) {
	conn, ctx := test.Setup(t)

	query := statement.Select("FirstName", "LastName").From("People").Where(conditional.Equal("FirstName", "Maddie")).Limit(1)

	var firstName, lastName string

	err := raptor.QueryRowStatement(ctx, conn, query).Scan(&firstName, &lastName)

	assert.NoError(t, err)

	assert.Equal(t, "Maddie", firstName)
	assert.Equal(t, "Schipper", lastName)
}

func TestQueryStatement(t *testing.T) {
	conn, ctx := test.Setup(t)

	query := statement.Select("FirstName", "LastName").From("People").Where(conditional.Equal("FirstName", "Maddie")).Limit(1)

	rows, err := raptor.QueryStatement(ctx, conn, query)
	require.NoError(t, err)

	assert.True(t, rows.Next())

	var firstName, lastName string
	err = rows.Scan(&firstName, &lastName)
	assert.NoError(t, err)
	assert.Equal(t, "Maddie", firstName)
	assert.Equal(t, "Schipper", lastName)

	assert.False(t, rows.Next())
}

func TestExecStatement(t *testing.T) {
	conn, ctx := test.Setup(t)

	exec := statement.Insert().Into("People").ValueMap(map[string]any{
		"FirstName": "Taylor",
		"LastName":  "Swift",
	})

	_, err := raptor.ExecStatement(ctx, conn, exec)
	assert.NoError(t, err)
}

func TestQueryStatementInsert(t *testing.T) {
	conn, ctx := test.Setup(t)

	exec := statement.Insert().Into("People").ValueMap(map[string]any{
		"FirstName": "Taylor",
		"LastName":  "Swift",
	}).Returning("FirstName", "LastName")

	var first, last string
	err := raptor.QueryRowStatement(ctx, conn, exec).Scan(&first, &last)

	assert.NoError(t, err)
	assert.Equal(t, "Taylor", first)
	assert.Equal(t, "Swift", last)
}

func TestQueryRowStatement_Error(t *testing.T) {
	row := raptor.QueryRowStatement(context.Background(), nil, new(raptortest.FailureGenerator))

	assert.Error(t, row.Err())
}

func TestQueryStatement_Error(t *testing.T) {
	_, err := raptor.QueryStatement(context.Background(), nil, new(raptortest.FailureGenerator))

	assert.Error(t, err)
}

func TestExecStatement_Error(t *testing.T) {
	_, err := raptor.ExecStatement(context.Background(), nil, new(raptortest.FailureGenerator))

	assert.Error(t, err)
}

func TestConn_QueryRowStatement_Error(t *testing.T) {
	conn, ctx := test.Setup(t)
	row := conn.QueryRowStatement(ctx, new(raptortest.FailureGenerator))

	assert.Error(t, row.Err())
}

func TestConn_QueryStatement_Error(t *testing.T) {
	conn, ctx := test.Setup(t)
	_, err := conn.QueryStatement(ctx, new(raptortest.FailureGenerator))

	assert.Error(t, err)
}

func TestConn_ExecStatement_Error(t *testing.T) {
	conn, ctx := test.Setup(t)
	_, err := conn.ExecStatement(ctx, new(raptortest.FailureGenerator))

	assert.Error(t, err)
}
