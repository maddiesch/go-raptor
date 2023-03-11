package raptor_test

import (
	"testing"

	"github.com/maddiesch/go-raptor"
	"github.com/maddiesch/go-raptor/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanRecord(t *testing.T) {
	conn, ctx := test.Setup(t)
	defer conn.Close()

	t.Run("Row", func(t *testing.T) {
		row := conn.QueryRow(ctx, `SELECT * FROM "People" WHERE "FirstName" = ? LIMIT 1;`, "Maddie")

		record, err := raptor.ScanRecord(row)
		require.NoError(t, err)

		assert.Equal(t, "Maddie", record["FirstName"])
		assert.Equal(t, "Schipper", record["LastName"])
	})

	t.Run("Rows", func(t *testing.T) {
		rows, err := conn.Query(ctx, `SELECT * FROM "People" WHERE "FirstName" = ? LIMIT 1;`, "Elle")
		require.NoError(t, err)

		for rows.Next() {
			record, err := raptor.ScanRecord(rows)
			require.NoError(t, err)

			assert.Equal(t, "Elle", record["FirstName"])
			assert.Equal(t, "Woods", record["LastName"])
		}

		require.NoError(t, rows.Err())
	})
}

func TestScanAllRecord(t *testing.T) {
	conn, ctx := test.Setup(t)
	defer conn.Close()

	rows, err := conn.Query(ctx, `SELECT * FROM "People";`)
	require.NoError(t, err)

	records, err := raptor.ScanAllRecord(rows)
	require.NoError(t, err)
	require.NoError(t, rows.Err())

	assert.Len(t, records, 3)
}

func TestRecordUnmarshal(t *testing.T) {
	conn, ctx := test.Setup(t)
	defer conn.Close()

	t.Run("Unmarshal Pet with all columns", func(t *testing.T) {
		var pet test.Pet

		row := conn.QueryRow(ctx, `SELECT * FROM "Pets" WHERE "Name" = ? LIMIT 1;`, "Sterling")
		err := raptor.UnmarshalRow(row, &pet)

		require.NoError(t, err)

		assert.Equal(t, "Sterling", pet.Name)

		if assert.NotNil(t, pet.Age) {
			assert.Equal(t, 5, *pet.Age)
		}
	})

	t.Run("Unmarshal Pet with select columns", func(t *testing.T) {
		var pet test.Pet

		row := conn.QueryRow(ctx, `SELECT "Name", "ParentID" FROM "Pets" WHERE "Name" = ? LIMIT 1;`, "Sterling")
		err := raptor.UnmarshalRow(row, &pet)

		require.NoError(t, err)

		assert.Equal(t, "Sterling", pet.Name)
		assert.NotEqual(t, 0, pet.PersonID)
	})

	t.Run("Unmarshal Pet nil age", func(t *testing.T) {
		var pet test.Pet

		row := conn.QueryRow(ctx, `SELECT * FROM "Pets" WHERE "Name" = ? LIMIT 1;`, "Lulu")
		err := raptor.UnmarshalRow(row, &pet)

		require.NoError(t, err)

		assert.Equal(t, "Lulu", pet.Name)
		assert.Nil(t, pet.Age)
	})

	t.Run("setting non-existing columns", func(t *testing.T) {
		var pet test.Pet

		row := conn.QueryRow(ctx, `SELECT "Name", 1 as "One" FROM "Pets" WHERE "Name" = ? LIMIT 1;`, "Lulu")
		err := raptor.UnmarshalRow(row, &pet)

		require.NoError(t, err)

		assert.Equal(t, "Lulu", pet.Name)
		assert.Nil(t, pet.Age)
	})
}
