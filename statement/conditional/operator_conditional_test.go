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
}

func TestNotEqual(t *testing.T) {
	pro := generator.NewIncrementingArgumentNameProvider()
	out, args := conditional.NotEqual("Test", "Foo").Generate(pro)

	assert.Equal(t, `"Test" != $v1`, out)
	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", "Foo"), args[0])
	}
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
