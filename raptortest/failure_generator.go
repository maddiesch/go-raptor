package raptortest

import (
	"errors"

	"github.com/maddiesch/go-raptor/statement/generator"
)

type FailureGenerator struct {
	Err error
}

func (f *FailureGenerator) Generate() (string, []any, error) {
	if f.Err != nil {
		return "", nil, f.Err
	}

	return "", nil, errors.New("FailureGenerator.Generate")
}

var (
	_ generator.Generator = (*FailureGenerator)(nil)
)
