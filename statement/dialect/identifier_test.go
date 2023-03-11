package dialect_test

import (
	"testing"

	"github.com/maddiesch/go-raptor/statement/dialect"
	"github.com/stretchr/testify/assert"
)

func TestIdentifier(t *testing.T) {
	assert.Equal(t, `"Foo"`, dialect.Identifier("Foo"))
}
