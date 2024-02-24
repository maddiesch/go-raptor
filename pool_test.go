package raptor_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	p := raptor.NewPool(5, func(context.Context) (*raptor.Conn, error) {
		return raptor.New(tempPath + "/testing.db?cache=shared")
	})
	t.Cleanup(func() {
		require.NoError(t, p.Close(context.Background()))
	})

	_, err := p.Exec(context.Background(), `CREATE TABLE "TestTable" ("ID" INTEGER NOT NULL PRIMARY KEY, "Index" INTEGER);`)
	require.NoError(t, err)

	waitGroup, ctx := errgroup.WithContext(context.Background())

	for v := 0; v < 100; v++ {
		v := v
		waitGroup.Go(func() error {
			_, err := p.Exec(ctx, `INSERT INTO "TestTable" ("Index") VALUES (?);`, v)
			assert.NoError(t, err)
			return nil
		})
	}

	for v := 0; v < 10; v++ {
		v := v
		waitGroup.Go(func() error {
			return p.Transact(ctx, func(conn raptor.DB) error {
				_, err := conn.Exec(ctx, `INSERT INTO "TestTable" ("Index") VALUES (?);`, v)

				assert.NoError(t, err)

				return nil
			})
		})
	}

	for v := 0; v < 100; v++ {
		waitGroup.Go(func() error {
			var count int64

			err := p.QueryRow(ctx, `SELECT COUNT(*) FROM "TestTable";`).Scan(&count)

			assert.NoError(t, err)

			return nil
		})
	}

	for v := 0; v < 8; v++ {
		v := v
		waitGroup.Go(func() error {
			var id int64

			err := p.ForWriting(ctx, func(d raptor.DB) error {
				return d.QueryRow(ctx, `INSERT INTO "TestTable" ("Index") VALUES (?) RETURNING rowid;`, v).Scan(&id)
			})

			assert.NoError(t, err)
			assert.NotEqual(t, 0, id)

			return nil
		})
	}

	waitGroup.Go(func() error {
		rows, err := p.Query(ctx, `SELECT COUNT(*) FROM "TestTable";`)
		require.NoError(t, err)
		defer rows.Close()

		for rows.Next() {

		}

		return nil
	})

	require.NoError(t, waitGroup.Wait())

	t.Run("Pool.QueryRow", func(t *testing.T) {
		p := raptor.NewPool(1, func(context.Context) (*raptor.Conn, error) {
			return raptor.New("file:testing.db?cache=shared&mode=memory")
		})
		t.Cleanup(func() {
			require.NoError(t, p.Close(context.Background()))
		})

		_, done, err := p.Reader(context.Background())
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, done())
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		t.Cleanup(cancel)

		row := p.QueryRow(ctx, `SELECT COUNT(*) FROM "TestTable";`)
		assert.ErrorIs(t, row.Err(), context.DeadlineExceeded)
		assert.ErrorIs(t, row.Scan(), context.DeadlineExceeded)

		_, err = row.Columns()
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
