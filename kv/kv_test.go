package kv_test

import (
	"context"
	"testing"

	"github.com/maddiesch/go-raptor/internal/test"
	"github.com/maddiesch/go-raptor/kv"
	"github.com/maddiesch/go-raptor/raptortest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepare(t *testing.T) {
	conn, ctx := test.Setup(t)

	err := kv.Prepare(ctx, conn)
	assert.NoError(t, err)

	t.Run("failure state", func(t *testing.T) {
		err := kv.Prepare(ctx, &raptortest.FailureConn{})
		assert.Error(t, err)
	})
}

func TestCRUD(t *testing.T) {
	conn, ctx := test.Setup(t)

	err := kv.Prepare(ctx, conn)
	require.NoError(t, err)

	err = kv.Set(ctx, conn, "test-key", []byte("value"))
	assert.NoError(t, err)

	err = kv.Set(ctx, conn, "test-key", []byte("value2"))
	assert.NoError(t, err)

	val, err := kv.Get(ctx, conn, "test-key")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), val)

	assert.True(t, kv.Exists(ctx, conn, "test-key"))

	err = kv.Delete(ctx, conn, "test-key")
	assert.NoError(t, err)

	assert.False(t, kv.Exists(ctx, conn, "test-key"))
}

func TestExists(t *testing.T) {
	conn := new(raptortest.FailureConn)

	v := kv.Exists(context.Background(), conn, "test-key")
	assert.False(t, v)
}
