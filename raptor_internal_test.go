package raptor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRowErrorChecking(t *testing.T) {
	row := &Row{err: errors.New("testing error")}

	t.Run("Scan", func(t *testing.T) {
		assert.Error(t, row.Scan())
	})

	t.Run("Columns", func(t *testing.T) {
		_, err := row.Columns()
		assert.Error(t, err)
	})

	t.Run("Err", func(t *testing.T) {
		assert.Error(t, row.Err())
	})
}

func Test_txConn(t *testing.T) {
	conn, err := New("file:internal-transaction?mode=memory&cache=shared")
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("multi-begin calls", func(t *testing.T) {
		conn.Transact(ctx, func(d DB) error {
			tx := d.(*txConn)

			err := tx.begin(ctx)
			assert.ErrorIs(t, err, ErrTransactionAlreadyStarted)

			return nil
		})
	})

	t.Run("multi-commit calls", func(t *testing.T) {
		conn.Transact(ctx, func(d DB) error {
			tx := d.(*txConn)

			err := tx.commit(ctx)
			assert.NoError(t, err)

			err = tx.commit(ctx)
			assert.NoError(t, err)

			return nil
		})
	})

	t.Run("multi-rollback calls", func(t *testing.T) {
		conn.Transact(ctx, func(d DB) error {
			tx := d.(*txConn)

			err := tx.rollback(ctx)
			assert.NoError(t, err)

			err = tx.rollback(ctx)
			assert.NoError(t, err)

			return nil
		})
	})
}
