//go:build go1.21

package raptor_test

import (
	"database/sql"
	"log/slog"
	"strings"
	"testing"

	"github.com/maddiesch/go-raptor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSLogFuncQueryLogger(t *testing.T) {
	var out strings.Builder

	conn, ctx := createTestConnection(t)
	defer conn.Close()

	logger := slog.New(slog.NewTextHandler(&out, nil))

	conn.SetLogger(
		raptor.NewSLogFuncQueryLogger(logger.InfoContext),
	)

	err := conn.Transact(ctx, func(tx raptor.DB) error {
		var v int64
		return tx.QueryRow(ctx, "SELECT $v1 AS one;", sql.Named("v1", 1)).Scan(&v)
	})
	require.NoError(t, err)

	assert.Contains(t, out.String(), "SAVEPOINT ")
	assert.Contains(t, out.String(), "RELEASE SAVEPOINT ")
}
