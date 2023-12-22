package raptor_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/maddiesch/go-raptor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestPool(t *testing.T) {
	tempPath := filepath.Join(t.TempDir(), "raptor_pool_test")
	require.NoError(t, os.MkdirAll(tempPath, os.ModeDir|0755))
	t.Cleanup(func() {
		os.RemoveAll(tempPath)
	})

	pool := raptor.NewPool(5, func(context.Context) (*raptor.Conn, error) {
		conn, err := raptor.New(tempPath + "/testing.db?cache=shared")
		if err != nil {
			return nil, err
		}

		return conn, nil
	})
	t.Cleanup(func() {
		require.NoError(t, pool.Close(context.Background()))
	})

	_, err := pool.Exec(context.Background(), `CREATE TABLE "TestTable" ("ID" INTEGER NOT NULL PRIMARY KEY, "Index" INTEGER);`)
	require.NoError(t, err)

	waitGroup, ctx := errgroup.WithContext(context.Background())

	for v := 0; v < 100; v++ {
		v := v
		waitGroup.Go(func() error {
			_, err := pool.Exec(ctx, `INSERT INTO "TestTable" ("Index") VALUES (?);`, v)
			assert.NoError(t, err)
			return nil
		})
	}

	for v := 0; v < 10; v++ {
		v := v
		waitGroup.Go(func() error {
			return pool.Transact(ctx, func(conn raptor.DB) error {
				_, err := conn.Exec(ctx, `INSERT INTO "TestTable" ("Index") VALUES (?);`, v)

				assert.NoError(t, err)

				return nil
			})
		})
	}

	for v := 0; v < 100; v++ {
		waitGroup.Go(func() error {
			var count int64

			err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM "TestTable";`).Scan(&count)

			assert.NoError(t, err)

			return nil
		})
	}

	require.NoError(t, waitGroup.Wait())
}
