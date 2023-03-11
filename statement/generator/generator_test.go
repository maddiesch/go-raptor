package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_incrementingArgumentNameProvider(t *testing.T) {
	t.Run("Next", func(t *testing.T) {
		p := NewIncrementingArgumentNameProvider()

		assert.Equal(t, "v1", p.Next())
		assert.Equal(t, "v2", p.Next())
	})
}
