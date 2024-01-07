package raptor_test

import (
	"testing"

	"github.com/maddiesch/go-raptor"
	"github.com/stretchr/testify/assert"
)

func TestTransactionHelperFunctions(t *testing.T) {
	conn, ctx := createTestConnection(t)

	t.Run("raptor.Transact", func(t *testing.T) {
		err := raptor.Transact(ctx, conn, func(d raptor.DB) error {
			return nil
		})

		assert.NoError(t, err)
	})

	t.Run("raptor.TransactV", func(t *testing.T) {
		count, err := raptor.TransactV(ctx, conn, func(d raptor.DB) (c int64, err error) {
			err = d.QueryRow(ctx, `SELECT COUNT(*) FROM "TestTable"`).Scan(&c)
			return
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})

	t.Run("raptor.TransactV2", func(t *testing.T) {
		name, age, err := raptor.TransactV2(ctx, conn, func(d raptor.DB) (n string, a int64, err error) {
			err = d.QueryRow(ctx, `SELECT "Name", "Age" FROM "TestTable" ORDER BY rowid ASC LIMIT 1`).Scan(&n, &a)
			return
		})

		assert.NoError(t, err)
		assert.Equal(t, "test", name)
		assert.Equal(t, int64(100), age)
	})
}
