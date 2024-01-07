package statement_test

import (
	"testing"

	"github.com/maddiesch/go-raptor/statement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTableBuilder(t *testing.T) {
	t.Run("basic query generation", func(t *testing.T) {
		s, _, err := statement.CreateTable("People").PrimaryKey("ID", statement.ColumnTypeText).Column(
			statement.Column("FirstName", statement.ColumnTypeText).NotNull(),
		).Generate()

		require.NoError(t, err)
		assert.Equal(t, `CREATE TABLE "People" ("ID" TEXT PRIMARY KEY NOT NULL UNIQUE, "FirstName" TEXT NOT NULL);`, s)
	})

	t.Run("if not exists", func(t *testing.T) {
		s, _, err := statement.CreateTable("People").IfNotExists().PrimaryKey("ID", statement.ColumnTypeText).Column(
			statement.Column("UserName", statement.ColumnTypeText).NotNull().Unique(),
			statement.Column("FirstName", statement.ColumnTypeText).NotNull(),
			statement.Column("CreatedAt", statement.ColumnTypeInteger).NotNull().Default("CURRENT_TIMESTAMP"),
		).Generate()

		require.NoError(t, err)
		assert.Equal(t, `CREATE TABLE IF NOT EXISTS "People" ("ID" TEXT PRIMARY KEY NOT NULL UNIQUE, "UserName" TEXT NOT NULL UNIQUE, "FirstName" TEXT NOT NULL, "CreatedAt" INTEGER NOT NULL DEFAULT CURRENT_TIMESTAMP);`, s)
	})
}
