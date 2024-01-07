package raptor

import (
	"context"
	"sync"

	"github.com/maddiesch/go-raptor/pool"
)

// Pool implements a thread-safe pool of database connections.
//
// It works as a single-writer multi-reader connection.
type Pool struct {
	pool.Pool[*Conn]

	wLock sync.RWMutex
}

// Create a new pool with the given number of maximum connections.
func NewPool(size int64, fn func(context.Context) (*Conn, error)) *Pool {
	config := pool.Config{
		MaxSize: size,
	}

	return &Pool{
		Pool: pool.New[*Conn](config, fn),
	}
}

func (p *Pool) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	return pool.WithValue(ctx, p.Pool, func(conn *Conn) (Result, error) {
		p.wLock.Lock()
		defer p.wLock.Unlock()

		return conn.Exec(ctx, query, args...)
	})
}

func (p *Pool) Query(ctx context.Context, query string, args ...any) (*Rows, error) {
	return pool.WithValue(ctx, p.Pool, func(conn *Conn) (*Rows, error) {
		p.wLock.RLock()
		defer p.wLock.RUnlock()

		return conn.Query(ctx, query, args...)
	})
}

func (p *Pool) QueryRow(ctx context.Context, query string, args ...any) *Row {
	row, err := pool.WithValue(ctx, p.Pool, func(conn *Conn) (*Row, error) {
		p.wLock.RLock()
		defer p.wLock.RUnlock()

		return conn.QueryRow(ctx, query, args...), nil
	})
	if err != nil {
		panic(err)
	}
	return row
}

func (p *Pool) Transact(ctx context.Context, fn func(DB) error) error {
	return pool.With(ctx, p.Pool, func(conn *Conn) error {
		p.wLock.Lock()
		defer p.wLock.Unlock()

		return conn.Transact(ctx, fn)
	})
}

var _ DB = (*Pool)(nil)
