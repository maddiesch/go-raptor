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

	t.Run("Get Timeout", func(t *testing.T) {
		v1, _ := p.Get(context.Background())
		v2, _ := p.Get(context.Background())

		t.Cleanup(func() {
			p.Put(v1)
			p.Put(v2)
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		t.Cleanup(cancel)
		_, err := p.Get(ctx)

		assert.ErrorIs(t, err, context.DeadlineExceeded)
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

func TestWithValue(t *testing.T) {
	var value atomic.Int64

	p := pool.New[int64](pool.Config{MaxSize: 2}, func(context.Context) (int64, error) {
		return value.Add(1), nil
	})

	t.Cleanup(func() {
		p.Close(context.Background())
	})

	result, err := pool.WithValue(context.Background(), p, func(s int64) (int64, error) {
		return s, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result)
}

func TestWithValue2(t *testing.T) {
	var value atomic.Int64

	p := pool.New[int64](pool.Config{MaxSize: 2}, func(context.Context) (int64, error) {
		return value.Add(1), nil
	})

	t.Cleanup(func() {
		p.Close(context.Background())
	})

	v1, v2, err := pool.WithValue2(context.Background(), p, func(s int64) (int64, int64, error) {
		return s, s + 1, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), v1)
	assert.Equal(t, int64(2), v2)
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
