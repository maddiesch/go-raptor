package raptortest_test

import (
	"errors"
	"testing"

	"github.com/maddiesch/go-raptor/raptortest"
	"github.com/stretchr/testify/assert"
)

func TestFailureGenerator(t *testing.T) {
	t.Run("DefaultError", func(t *testing.T) {
		_, _, err := new(raptortest.FailureGenerator).Generate()

		assert.Error(t, err)
	})

	t.Run("CustomError", func(t *testing.T) {
		g := &raptortest.FailureGenerator{}
		g.Err = errors.New("custom error")

		_, _, err := g.Generate()

		assert.ErrorIs(t, err, g.Err)
	})
}
