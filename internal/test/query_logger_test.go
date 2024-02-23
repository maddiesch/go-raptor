package test_test

import (
	"testing"

	"github.com/maddiesch/go-raptor/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestCollectQueryLogger(t *testing.T) {
	conn, ctx := test.Setup(t)
	log := new(test.CollectQueryLogger)

	conn.SetQueryLogger(log)

	_, err := conn.Exec(ctx, "SELECT 1")
	assert.NoError(t, err)
	assert.Len(t, log.Queries, 1)
	assert.Equal(t, "SELECT 1", log.Queries[0].Query)

	log.Reset()
	assert.Len(t, log.Queries, 0)
}

func TestTestQueryLogger(t *testing.T) {
	conn, ctx := test.Setup(t)
	conn.SetLogger(&test.TestQueryLogger{TestingT: t})

	assert.NotPanics(t, func() {
		_, err := conn.Exec(ctx, "SELECT 1")
		assert.NoError(t, err)
	})
}
