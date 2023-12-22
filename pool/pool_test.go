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
}
