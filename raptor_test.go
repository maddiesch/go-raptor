package raptor_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/maddiesch/go-raptor"
	"github.com/maddiesch/go-raptor/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestConnection(t test.TestingT) (*raptor.Conn, context.Context) {
	nh := md5.Sum([]byte(t.Name()))
	ctx := context.Background()
	conn, err := raptor.New(fmt.Sprintf("file:%s?mode=memory&cache=shared", hex.EncodeToString(nh[:])))
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `CREATE TABLE "TestTable" ("ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "Name" TEXT NOT NULL DEFAULT '', "Age" INTEGER NOT NULL DEFAULT 0);`)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `INSERT INTO "TestTable" ("Name", "Age") VALUES ('test', 100);`)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `INSERT INTO "TestTable" ("Name", "Age") VALUES ('test-two', 200);`)
	require.NoError(t, err)

	conn.SetLogger(&test.TestQueryLogger{TestingT: t})

	return conn, ctx
}

func TestNewConn(t *testing.T) {
	conn, err := raptor.New("file:test-new-conn?mode=memory&cache=shared")
	require.NoError(t, err)

	assert.NoError(t, conn.Close())
}

func TestConn_Ping(t *testing.T) {
	conn, ctx := createTestConnection(t)
	defer conn.Close()

	assert.NoError(t, conn.Ping(ctx))
}

func TestConn_Transact(t *testing.T) {
	conn, ctx := createTestConnection(t)
	defer conn.Close()

	t.Run("insert", func(t *testing.T) {
		err := conn.Transact(ctx, func(tx raptor.DB) error {
			_, err := tx.Exec(ctx, `INSERT INTO "TestTable" DEFAULT VALUES;`)
			return err
		})

		assert.NoError(t, err)
	})

	t.Run("nested", func(t *testing.T) {
		err := conn.Transact(ctx, func(tx raptor.DB) error {
			_, err := tx.Exec(ctx, `INSERT INTO "TestTable" DEFAULT VALUES;`)
			require.NoError(t, err)

			return tx.Transact(ctx, func(tx raptor.DB) error {
				_, err := tx.Exec(ctx, `INSERT INTO "TestTable" DEFAULT VALUES;`)
				return err
			})
		})

		assert.NoError(t, err)
	})

	t.Run("rollback when the function panics then re-panic", func(t *testing.T) {
		conn, ctx := createTestConnection(t)
		defer conn.Close()

		coll := &test.CollectQueryLogger{}

		conn.SetLogger(coll)

		assert.Panics(t, func() {
			conn.Transact(ctx, func(raptor.DB) error {
				panic("expected to panic")
			})
		})

		if assert.Len(t, coll.Queries, 2) {
			assert.Contains(t, coll.Queries[0].Query, "SAVEPOINT ")
			assert.Contains(t, coll.Queries[1].Query, "ROLLBACK TRANSACTION TO SAVEPOINT ")
		}
	})

	t.Run("rollback happens when the transaction function returns an error", func(t *testing.T) {
		conn, ctx := createTestConnection(t)
		defer conn.Close()

		coll := &test.CollectQueryLogger{}

		conn.SetLogger(coll)

		eRollback := errors.New("trigger rollback")

		err := conn.Transact(ctx, func(raptor.DB) error {
			return eRollback
		})

		require.ErrorIs(t, err, eRollback)

		if assert.Len(t, coll.Queries, 2) {
			assert.Contains(t, coll.Queries[0].Query, "SAVEPOINT ")
			assert.Contains(t, coll.Queries[1].Query, "ROLLBACK TRANSACTION TO SAVEPOINT ")
		}
	})

	t.Run("force a rollback with ErrTxRollback", func(t *testing.T) {
		conn, ctx := createTestConnection(t)
		defer conn.Close()

		err := conn.Transact(ctx, func(raptor.DB) error {
			return raptor.ErrTxRollback
		})

		assert.NoError(t, err)
	})

	t.Run("query inside of a transaction returns data", func(t *testing.T) {
		conn, ctx := createTestConnection(t)
		defer conn.Close()

		err := conn.Transact(ctx, func(tx raptor.DB) error {
			if _, err := tx.Exec(ctx, `INSERT INTO "TestTable" DEFAULT VALUES;`); err != nil {
				return err
			}

			var count int
			if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM "TestTable";`).Scan(&count); err != nil {
				return err
			}

			assert.Equal(t, 3, count)

			return raptor.ErrTxRollback
		})

		require.NoError(t, err)

		var count int
		err = conn.QueryRow(ctx, `SELECT COUNT(*) FROM "TestTable";`).Scan(&count)
		require.NoError(t, err)

		assert.Equal(t, 2, count)
	})

	t.Run("using a transaction outside of the block returns an error", func(t *testing.T) {
		conn, ctx := createTestConnection(t)
		defer conn.Close()

		var tx raptor.DB

		err := conn.Transact(ctx, func(itx raptor.DB) error {
			tx = itx
			return nil
		})

		require.NoError(t, err)
		require.NotNil(t, tx)

		_, err = tx.Exec(ctx, `INSERT INTO "TestTable" DEFAULT VALUES:`)
		assert.ErrorIs(t, err, raptor.ErrTransactionNotRunning)

		err = tx.Transact(ctx, func(raptor.DB) error { return nil })
		assert.ErrorIs(t, err, raptor.ErrTransactionNotRunning)

		_, err = tx.Query(ctx, `SELECT * FROM "TestTable";`)
		assert.ErrorIs(t, err, raptor.ErrTransactionNotRunning)

		row := tx.QueryRow(ctx, `SELECT * FROM "TestTable";`)
		assert.ErrorIs(t, row.Err(), raptor.ErrTransactionNotRunning)
	})
}

func TestConn_Query(t *testing.T) {
	t.Run("query for a single row", func(t *testing.T) {
		conn, ctx := createTestConnection(t)
		defer conn.Close()

		rows, err := conn.Query(ctx, `SELECT * FROM "TestTable" WHERE "Name" = ? LIMIT 1;`, "test")
		require.NoError(t, err)

		var count int
		for rows.Next() {
			count++

			var name string
			var id, age int

			err := rows.Scan(&id, &name, &age)

			require.NoError(t, err)

			assert.Equal(t, "test", name)
			assert.Equal(t, 100, age)
		}

		assert.Equal(t, 1, count)
	})
}

func TestConn_QueryRow(t *testing.T) {
	t.Run("query for a single row", func(t *testing.T) {
		conn, ctx := createTestConnection(t)
		defer conn.Close()

		var name string
		var id, age int
		err := conn.QueryRow(ctx, `SELECT * FROM "TestTable" WHERE "Name" = ? LIMIT 1;`, "test").Scan(&id, &name, &age)
		require.NoError(t, err)

		assert.Equal(t, "test", name)
		assert.Equal(t, 100, age)
	})

	t.Run("returns expected error for no results", func(t *testing.T) {
		conn, ctx := createTestConnection(t)
		defer conn.Close()

		err := conn.QueryRow(ctx, `SELECT * FROM "TestTable" WHERE 1 = 2;`).Scan()

		assert.ErrorIs(t, err, raptor.ErrNoRows)
	})
}

func TestNewQueryLogger(t *testing.T) {
	var out strings.Builder

	conn, ctx := createTestConnection(t)
	defer conn.Close()

	conn.SetLogger(raptor.NewQueryLogger(&out))

	err := conn.Transact(ctx, func(raptor.DB) error { return nil })
	require.NoError(t, err)

	assert.Contains(t, out.String(), "SAVEPOINT ")
	assert.Contains(t, out.String(), "RELEASE SAVEPOINT ")
}
