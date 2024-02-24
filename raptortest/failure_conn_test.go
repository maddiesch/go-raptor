package raptortest_test

import (
	"context"
	"testing"

	"github.com/maddiesch/go-raptor"
	"github.com/maddiesch/go-raptor/raptortest"
	"github.com/stretchr/testify/assert"
)

func TestFailureConn(t *testing.T) {
	c := new(raptortest.FailureConn)

	t.Run("Exec", func(t *testing.T) {
		_, err := c.Exec(context.Background(), "SELECT 1")
		assert.Error(t, err)
	})

	t.Run("Query", func(t *testing.T) {
		_, err := c.Query(context.Background(), "SELECT 1")
		assert.Error(t, err)
	})

	t.Run("QueryRow", func(t *testing.T) {
		row := c.QueryRow(context.Background(), "SELECT 1")
		assert.Error(t, row.Scan())
		assert.Error(t, row.Err())

		_, err := row.Columns()
		assert.Error(t, err)
	})

	t.Run("Transact", func(t *testing.T) {
		err := c.Transact(context.Background(), func(tx raptor.DB) error {
			_, err := tx.Exec(context.Background(), "SELECT 1")
			return err
		})
		assert.Error(t, err)
	})
}
