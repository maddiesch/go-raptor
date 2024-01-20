package conditional_test

import (
	"database/sql"
	"testing"

	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/maddiesch/go-raptor/statement/generator"
	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	pro := generator.NewIncrementingArgumentNameProvider()
	out, args := conditional.Equal("Test", "Foo").Generate(pro)

	assert.Equal(t, `"Test" = $v1`, out)
	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", "Foo"), args[0])
	}

	t.Run("when value is nil", func(t *testing.T) {
		out, args = conditional.Equal("Test", nil).Generate(pro)

		assert.Equal(t, `"Test" IS NULL`, out)
	})
}

func TestNotEqual(t *testing.T) {
	pro := generator.NewIncrementingArgumentNameProvider()
	out, args := conditional.NotEqual("Test", "Foo").Generate(pro)

	assert.Equal(t, `"Test" != $v1`, out)
	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", "Foo"), args[0])
	}

	t.Run("when value is nil", func(t *testing.T) {
		out, args = conditional.NotEqual("Test", nil).Generate(pro)

		assert.Equal(t, `"Test" IS NOT NULL`, out)
	})
}

func TestGreaterThan(t *testing.T) {
	pro := generator.NewIncrementingArgumentNameProvider()
	out, args := conditional.GreaterThan("Test", "Foo").Generate(pro)

	assert.Equal(t, `"Test" > $v1`, out)
	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", "Foo"), args[0])
	}
}

func TestGreaterThanEq(t *testing.T) {
	pro := generator.NewIncrementingArgumentNameProvider()
	out, args := conditional.GreaterThanEq("Test", "Foo").Generate(pro)

	assert.Equal(t, `"Test" >= $v1`, out)
	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", "Foo"), args[0])
	}
}

func TestLessThan(t *testing.T) {
	pro := generator.NewIncrementingArgumentNameProvider()
	out, args := conditional.LessThan("Test", "Foo").Generate(pro)

	assert.Equal(t, `"Test" < $v1`, out)
	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", "Foo"), args[0])
	}
}

func TestLessThanEq(t *testing.T) {
	pro := generator.NewIncrementingArgumentNameProvider()
	out, args := conditional.LessThanEq("Test", "Foo").Generate(pro)

	assert.Equal(t, `"Test" <= $v1`, out)
	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", "Foo"), args[0])
	}
}

func TestNull(t *testing.T) {
	pro := generator.NewIncrementingArgumentNameProvider()
	out, args := conditional.Null("Test").Generate(pro)

	assert.Equal(t, `"Test" IS NULL`, out)
	assert.Len(t, args, 0)
}

func TestNotNull(t *testing.T) {
	pro := generator.NewIncrementingArgumentNameProvider()
	out, args := conditional.NotNull("Test").Generate(pro)

	assert.Equal(t, `"Test" IS NOT NULL`, out)
	assert.Len(t, args, 0)
}
