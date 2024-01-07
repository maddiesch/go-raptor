//go:build go1.21

package raptor_test

import (
	"database/sql"
	"log/slog"
	"strings"
	"testing"

	"github.com/maddiesch/go-raptor"
	"github.com/maddiesch/go-raptor/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSLogFuncQueryLogger(t *testing.T) {
	t.Run("given named argument", func(t *testing.T) {
		var out strings.Builder

		conn, ctx := test.Setup(t)

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
	})

	t.Run("given literal arguments", func(t *testing.T) {
		var out strings.Builder

		conn, ctx := test.Setup(t)

		logger := slog.New(slog.NewTextHandler(&out, nil))

		conn.SetLogger(
			raptor.NewSLogFuncQueryLogger(logger.InfoContext),
		)

		err := conn.Transact(ctx, func(tx raptor.DB) error {
			var v int64
			var v2 string
			return tx.QueryRow(ctx, "SELECT ? AS one, ? AS two;", 1, "foo").Scan(&v, &v2)
		})
		require.NoError(t, err)

		assert.Contains(t, out.String(), "SAVEPOINT ")
		assert.Contains(t, out.String(), "RELEASE SAVEPOINT ")
	})
}

func TestCreateSlogArg(t *testing.T) {
	assert.Equal(t, slog.String("0", "foo"), raptor.CreateSlogArg(0, "foo"))
	assert.Equal(t, slog.Int64("0", 1), raptor.CreateSlogArg(0, 1))
	assert.Equal(t, slog.Int64("0", 1), raptor.CreateSlogArg(0, int64(1)))
	assert.Equal(t, slog.Uint64("0", 1), raptor.CreateSlogArg(0, uint(1)))
	assert.Equal(t, slog.Uint64("0", 1), raptor.CreateSlogArg(0, uint64(1)))
	assert.Equal(t, slog.Any("v1", "FooBar"), raptor.CreateSlogArg(0, sql.Named("v1", "FooBar")))
}
