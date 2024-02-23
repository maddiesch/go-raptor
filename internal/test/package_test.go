package test_test

import (
	"testing"

	"github.com/maddiesch/go-raptor/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	vPtr := test.Ptr(42)

	assert.Equal(t, 42, *vPtr)
}
