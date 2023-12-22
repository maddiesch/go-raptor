package pool

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Pool[T any] interface {
	Get(context.Context) (T, error)

	Put(T) error

	Close(context.Context) error

	Len() int
}

type Config struct {
	MaxSize int64
}

func New[T any](c Config, fn func(context.Context) (T, error)) Pool[T] {
	return &pool[T]{
		max:       c.MaxSize,
		builder:   fn,
		semaphore: semaphore.NewWeighted(c.MaxSize),
		values:    make([]T, 0, c.MaxSize),
	}
}

type pool[T any] struct {
	max       int64
	builder   func(context.Context) (T, error)
	semaphore *semaphore.Weighted
	mu        sync.Mutex
	values    []T
}

func (p *pool[T]) Get(ctx context.Context) (T, error) {
	if err := p.semaphore.Acquire(ctx, 1); err != nil {
		var v T
		return v, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.values) == 0 {
		return p.builder(ctx)
	}

	v := p.values[len(p.values)-1]
	p.values = p.values[:len(p.values)-1]

	return v, nil
}

func (p *pool[T]) Put(v T) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.semaphore.Release(1)

	p.values = append(p.values, v)

	return nil
}

func (p *pool[T]) Close(ctx context.Context) error {
	if err := p.semaphore.Acquire(ctx, p.max); err != nil {
		return err
	}
	defer p.semaphore.Release(p.max)

	p.mu.Lock()
	defer p.mu.Unlock()

	for _, v := range p.values {
		switch v := any(v).(type) {
		case Closer:
			v.Close()
		case CloseErr:
			if err := v.Close(); err != nil {
				return err
			}
		case CloseContext:
			v.Close(ctx)
		case CloseContextErr:
			if err := v.Close(ctx); err != nil {
				return err
			}
		}
	}

	p.values = make([]T, 0, p.max)

	return nil
}

func (p *pool[T]) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.values)
}

func With[T any](ctx context.Context, pool Pool[T], fn func(T) error) error {
	v, err := pool.Get(ctx)
	if err != nil {
		return err
	}
	defer pool.Put(v)

	return fn(v)
}

func WithValue[T any, V any](ctx context.Context, pool Pool[T], fn func(T) (V, error)) (V, error) {
	v, err := pool.Get(ctx)
	if err != nil {
		var v V
		return v, err
	}
	defer pool.Put(v)

	return fn(v)
}

func WithValue2[T any, V any, V2 any](ctx context.Context, pool Pool[T], fn func(T) (V, V2, error)) (V, V2, error) {
	v, err := pool.Get(ctx)
	if err != nil {
		var v V
		var vv V2
		return v, vv, err
	}
	defer pool.Put(v)

	return fn(v)
}

type Closer interface {
	Close()
}

type CloseErr interface {
	Close() error
}

type CloseContext interface {
	Close(context.Context)
}

type CloseContextErr interface {
	Close(context.Context) error
}
