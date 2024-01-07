package migrate_test

import (
	"context"
	"os"
	"testing"

	"github.com/maddiesch/go-raptor"
	"github.com/maddiesch/go-raptor/migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUp(t *testing.T) {
	t.Run("given a migration that has not been run", func(t *testing.T) {
		db, err := raptor.New(":memory:?mode=memory&cache=shared")
		require.NoError(t, err)
		t.Cleanup(func() {
			db.Close()
		})

		db.SetLogger(raptor.NewQueryLogger(os.Stderr))

		m := migrate.Migration{
			Name: "testing-migration",
			Up: []string{
				`CREATE TABLE "migration_table" ("example" TEXT);`,
			},
		}

		err = migrate.Up(context.Background(), db, m)
		require.NoError(t, err)
	})

	t.Run("given a migration that has already been executed", func(t *testing.T) {
		db, err := raptor.New(":memory:?mode=memory&cache=shared")
		require.NoError(t, err)
		t.Cleanup(func() {
			db.Close()
		})

		db.SetLogger(raptor.NewQueryLogger(os.Stderr))

		m := migrate.Migration{
			Name: "testing-migration-dup",
			Up: []string{
				`CREATE TABLE "migration_table" ("example" TEXT);`,
			},
		}

		err = migrate.Up(context.Background(), db, m)
		require.NoError(t, err)

		err = migrate.Up(context.Background(), db, m)
		require.NoError(t, err)
	})
}

func TestDown(t *testing.T) {
	db, err := raptor.New(":memory:?mode=memory&cache=shared")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close()
	})

	db.SetLogger(raptor.NewQueryLogger(os.Stderr))

	m := []migrate.Migration{
		{
			Name: "testing-migration",
			Up: []string{
				`CREATE TABLE "migration_table" ("example" TEXT);`,
			},
			Down: []string{
				`DROP TABLE "migration_table";`,
			},
		},
	}

	err = migrate.Up(context.Background(), db, m...)
	require.NoError(t, err)

	m = append(m, migrate.Migration{
		Name: "testing-down-not-run",
	})

	err = migrate.Down(context.Background(), db, m...)
	assert.NoError(t, err)
}
