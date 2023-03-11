package generator

import (
	"fmt"
	"sync"
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
	mu      sync.Mutex
	current uint64
}

func (i *incrementingArgumentNameProvider) Next() string {
	i.mu.Lock()
	i.current++
	value := i.current
	i.mu.Unlock()

	return fmt.Sprintf("v%d", value)
}
