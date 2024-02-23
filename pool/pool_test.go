package pool_test

import (
	"context"
	"errors"
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

		t.Run("when error is returned", func(t *testing.T) {
			c := &poolValueCloserErr{err: errors.New(t.Name())}

			p := pool.New[any](pool.Config{MaxSize: 1}, func(ctx context.Context) (any, error) {
				return c, nil
			})

			pool.Load(context.Background(), p, 1)

			err := p.Close(context.Background())
			assert.ErrorIs(t, err, c.err)

			assert.True(t, c.called.Load())
		})
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

		t.Run("when error is returned", func(t *testing.T) {
			c := &poolValueCloserContextErr{err: errors.New(t.Name())}

			p := pool.New[any](pool.Config{MaxSize: 1}, func(ctx context.Context) (any, error) {
				return c, nil
			})

			pool.Load(context.Background(), p, 1)

			err := p.Close(context.Background())
			assert.ErrorIs(t, err, c.err)

			assert.True(t, c.called.Load())
		})
	})

	t.Run("Over Put", func(t *testing.T) {
		p := pool.New(pool.Config{MaxSize: 1}, func(_ context.Context) (int64, error) {
			return 0, nil
		})

		pool.Load(context.Background(), p, 1)

		assert.Panics(t, func() {
			p.Put(int64(1))
		})
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

type poolValueCloserErr struct {
	called atomic.Bool
	err    error
}

func (p *poolValueCloserErr) Close() error {
	p.called.Store(true)

	return p.err
}

type poolValueCloserContext struct {
	called atomic.Bool
}

func (p *poolValueCloserContext) Close(context.Context) {
	p.called.Store(true)
}

type poolValueCloserContextErr struct {
	called atomic.Bool
	err    error
}

func (p *poolValueCloserContextErr) Close(context.Context) error {
	p.called.Store(true)

	return p.err
}

var (
	_ pool.Closer          = (*poolValueCloser)(nil)
	_ pool.CloseErr        = (*poolValueCloserErr)(nil)
	_ pool.CloseContext    = (*poolValueCloserContext)(nil)
	_ pool.CloseContextErr = (*poolValueCloserContextErr)(nil)
)
