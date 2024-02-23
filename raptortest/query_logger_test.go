package raptortest_test

import (
	"testing"

	"github.com/maddiesch/go-raptor/internal/test"
	"github.com/maddiesch/go-raptor/raptortest"
	"github.com/stretchr/testify/assert"
)

func TestNewQueryLogger(t *testing.T) {
	l := raptortest.NewQueryLogger(t)
	db, ctx := test.Setup(t)

	db.SetLogger(l)

	_, err := db.Exec(ctx, `SELECT 1 as "One";`)

	assert.NoError(t, err)
}
