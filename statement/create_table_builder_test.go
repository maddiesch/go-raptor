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
}
