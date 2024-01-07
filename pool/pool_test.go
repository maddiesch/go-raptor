package pool_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maddiesch/go-raptor/pool"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestPool(t *testing.T) {
	var value atomic.Int64

	p := pool.New[int64](pool.Config{MaxSize: 2}, func(context.Context) (int64, error) {
		return value.Add(1), nil
	})

	t.Cleanup(func() {
		p.Close(context.Background())
	})

	t.Run("Get", func(t *testing.T) {
		group, ctx := errgroup.WithContext(context.Background())

		group.Go(func() error {
			return pool.With(ctx, p, func(s int64) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			})
		})

		group.Go(func() error {
			return pool.With(ctx, p, func(s int64) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			})
		})

		group.Go(func() error {
			return pool.With(ctx, p, func(s int64) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			})
		})

		group.Go(func() error {
			return pool.With(ctx, p, func(s int64) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			})
		})

		group.Go(func() error {
			return pool.With(ctx, p, func(s int64) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			})
		})

		assert.NoError(t, group.Wait())
	})

	t.Run("Closer", func(t *testing.T) {
		c := &poolValueCloser{}

		p := pool.New[any](pool.Config{MaxSize: 1}, func(ctx context.Context) (any, error) {
			return c, nil
		})

		pool.Load(context.Background(), p, 1)

		p.Close(context.Background())

		assert.True(t, c.called.Load())
	})

	t.Run("CloseErr", func(t *testing.T) {
		c := &poolValueCloserErr{}

		p := pool.New[any](pool.Config{MaxSize: 1}, func(ctx context.Context) (any, error) {
			return c, nil
		})

		pool.Load(context.Background(), p, 1)

		p.Close(context.Background())

		assert.True(t, c.called.Load())
	})

	t.Run("CloseContext", func(t *testing.T) {
		c := &poolValueCloserContext{}

		p := pool.New[any](pool.Config{MaxSize: 1}, func(ctx context.Context) (any, error) {
			return c, nil
		})

		pool.Load(context.Background(), p, 1)

		p.Close(context.Background())

		assert.True(t, c.called.Load())
	})

	t.Run("CloseContextErr", func(t *testing.T) {
		c := &poolValueCloserContextErr{}

		p := pool.New[any](pool.Config{MaxSize: 1}, func(ctx context.Context) (any, error) {
			return c, nil
		})

		pool.Load(context.Background(), p, 1)

		p.Close(context.Background())

		assert.True(t, c.called.Load())
	})
}

func TestLoad(t *testing.T) {
	var value atomic.Int64

	p := pool.New[int64](pool.Config{MaxSize: 2}, func(context.Context) (int64, error) {
		return value.Add(1), nil
	})

	t.Cleanup(func() {
		p.Close(context.Background())
	})

	err := pool.Load(context.Background(), p, 2)

	assert.NoError(t, err)
	assert.Equal(t, 2, p.Len())
}

type poolValueCloser struct {
	called atomic.Bool
}

func (p *poolValueCloser) Close() {
	p.called.Store(true)
}

var _ pool.Closer = (*poolValueCloser)(nil)

type poolValueCloserErr struct {
	called atomic.Bool
}

func (p *poolValueCloserErr) Close() error {
	p.called.Store(true)

	return nil
}

var _ pool.CloseErr = (*poolValueCloserErr)(nil)

type poolValueCloserContext struct {
	called atomic.Bool
}

func (p *poolValueCloserContext) Close(context.Context) {
	p.called.Store(true)
}

var _ pool.CloseContext = (*poolValueCloserContext)(nil)

type poolValueCloserContextErr struct {
	called atomic.Bool
}

func (p *poolValueCloserContextErr) Close(context.Context) error {
	p.called.Store(true)

	return nil
}

var _ pool.CloseContextErr = (*poolValueCloserContextErr)(nil)
