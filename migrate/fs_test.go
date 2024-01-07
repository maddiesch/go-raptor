package migrate_test

import (
	"embed"
	"testing"

	"github.com/maddiesch/go-raptor/migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testing/*.sql
var migrationFS embed.FS

func TestFromFS(t *testing.T) {
	migrations, err := migrate.FromFS(migrationFS)
	require.NoError(t, err)

	require.Len(t, migrations, 2)

	assert.Equal(t, "testing/1_setup", migrations[0].Name)
	assert.Len(t, migrations[0].Up, 1)
	assert.Len(t, migrations[0].Down, 1)

	assert.Equal(t, "testing/2_create_key_value", migrations[1].Name)
	assert.Len(t, migrations[1].Up, 1)
	assert.Len(t, migrations[1].Down, 0)
}
