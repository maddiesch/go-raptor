package generator

import (
	"fmt"
	"sync/atomic"
)

type Generator interface {
	Generate() (string, []any, error)
}

type ArgumentNameProvider interface {
	Next() string
}

func NewIncrementingArgumentNameProvider() ArgumentNameProvider {
	return &incrementingArgumentNameProvider{}
}

type incrementingArgumentNameProvider struct {
	v atomic.Int64
}

func (i *incrementingArgumentNameProvider) Next() string {
	value := i.v.Add(1)

	return fmt.Sprintf("v%d", value)
}
