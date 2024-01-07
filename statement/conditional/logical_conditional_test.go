package conditional_test

import (
	"database/sql"
	"testing"

	"github.com/maddiesch/go-raptor/statement/conditional"
	"github.com/maddiesch/go-raptor/statement/generator"
	"github.com/stretchr/testify/assert"
)

func TestConditionalAnd(t *testing.T) {
	t.Run("simple equality", func(t *testing.T) {
		provider := generator.NewIncrementingArgumentNameProvider()

		stmt, args := conditional.And(
			conditional.Equal("First", 1),
			conditional.Equal("Second", 2),
		).Generate(provider)

		assert.Equal(t, `("First" = $v1 AND "Second" = $v2)`, stmt)

		if assert.Len(t, args, 2) {
			assert.Equal(t, sql.Named("v1", 1), args[0])
			assert.Equal(t, sql.Named("v2", 2), args[1])
		}
	})

	t.Run("nested equality", func(t *testing.T) {
		provider := generator.NewIncrementingArgumentNameProvider()

		stmt, args := conditional.And(
			conditional.Equal("First", 1),
			conditional.And(
				conditional.Equal("Second", 2),
				conditional.Equal("Third", 3),
			),
		).Generate(provider)

		assert.Equal(t, `("First" = $v1 AND ("Second" = $v2 AND "Third" = $v3))`, stmt)

		if assert.Len(t, args, 3) {
			assert.Equal(t, sql.Named("v1", 1), args[0])
			assert.Equal(t, sql.Named("v2", 2), args[1])
			assert.Equal(t, sql.Named("v3", 3), args[2])
		}
	})

	t.Run("given nil lhs", func(t *testing.T) {
		provider := generator.NewIncrementingArgumentNameProvider()

		stmt, args := conditional.And(
			nil,
			conditional.Equal("RHS_Column", 2),
		).Generate(provider)

		assert.Equal(t, `"RHS_Column" = $v1`, stmt)

		if assert.Len(t, args, 1) {
			assert.Equal(t, sql.Named("v1", 2), args[0])
		}
	})

	t.Run("given nil rhs", func(t *testing.T) {
		provider := generator.NewIncrementingArgumentNameProvider()

		stmt, args := conditional.And(
			conditional.Equal("LHS_Column", 1),
			nil,
		).Generate(provider)

		assert.Equal(t, `"LHS_Column" = $v1`, stmt)

		if assert.Len(t, args, 1) {
			assert.Equal(t, sql.Named("v1", 1), args[0])
		}
	})
}

func TestConditionalOr(t *testing.T) {
	provider := generator.NewIncrementingArgumentNameProvider()

	stmt, args := conditional.Or(
		conditional.Equal("First", 1),
		conditional.Equal("Second", 2),
	).Generate(provider)

	assert.Equal(t, `("First" = $v1 OR "Second" = $v2)`, stmt)

	if assert.Len(t, args, 2) {
		assert.Equal(t, sql.Named("v1", 1), args[0])
		assert.Equal(t, sql.Named("v2", 2), args[1])
	}

	t.Run("when provided with no arguments", func(t *testing.T) {
		assert.Panics(t, func() {
			conditional.Or(nil, nil).Generate(provider)
		})
	})

	t.Run("given nil lhs", func(t *testing.T) {
		provider := generator.NewIncrementingArgumentNameProvider()

		stmt, args := conditional.Or(
			nil,
			conditional.Equal("RHS_Column", 2),
		).Generate(provider)

		assert.Equal(t, `"RHS_Column" = $v1`, stmt)

		if assert.Len(t, args, 1) {
			assert.Equal(t, sql.Named("v1", 2), args[0])
		}
	})

	t.Run("given nil rhs", func(t *testing.T) {
		provider := generator.NewIncrementingArgumentNameProvider()

		stmt, args := conditional.Or(
			conditional.Equal("LHS_Column", 1),
			nil,
		).Generate(provider)

		assert.Equal(t, `"LHS_Column" = $v1`, stmt)

		if assert.Len(t, args, 1) {
			assert.Equal(t, sql.Named("v1", 1), args[0])
		}
	})
}

func TestConditionalIn(t *testing.T) {
	provider := generator.NewIncrementingArgumentNameProvider()

	stmt, args := conditional.In("First", []int{1, 2, 3}).Generate(provider)

	assert.Equal(t, `("First" IN $v1)`, stmt)
	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", []int{1, 2, 3}), args[0])
	}
}

func TestConditionalNotIn(t *testing.T) {
	provider := generator.NewIncrementingArgumentNameProvider()

	stmt, args := conditional.NotIn("First", []int{1, 2, 3}).Generate(provider)

	assert.Equal(t, `("First" NOT IN $v1)`, stmt)
	if assert.Len(t, args, 1) {
		assert.Equal(t, sql.Named("v1", []int{1, 2, 3}), args[0])
	}
}
