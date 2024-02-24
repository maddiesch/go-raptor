package raptor

import (
	"context"
	"sync"

	"github.com/maddiesch/go-raptor/pool"
)

// Pool implements a thread-safe pool of database connections.
//
// It works as a single-writer multi-reader connection.
//
// Warning: The pool is not safe to perform mutating queries:
// pool.Query(ctx, "INSERT INTO table (id) VALUES (?) RETURNING id", ...)
// will not lock the pool for writing even though it's a mutating query.
//
// If you plan on using mutating queries, you should use ForWriting, or Transact
type Pool struct {
	pool.Pool[*Conn]

	wLock sync.RWMutex
}

// Create a new pool with the given number of maximum connections.
//
// Panic if size is less than 1.
func NewPool(size int64, fn func(context.Context) (*Conn, error)) *Pool {
	if size < 1 {
		panic("raptor: pool size must be at least 1")
	}

	config := pool.Config{
		MaxSize: size,
	}

	return &Pool{
		Pool: pool.New[*Conn](config, func(ctx context.Context) (*Conn, error) {
			return fn(ctx)
		}),
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

func (p *Pool) QueryRow(ctx context.Context, query string, args ...any) Row {
	row, err := pool.WithValue(ctx, p.Pool, func(conn *Conn) (Row, error) {
		p.wLock.RLock()
		defer p.wLock.RUnlock()

		return conn.QueryRow(ctx, query, args...), nil
	})
	if err != nil {
		return &poolRowErr{err}
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

// ForWriting is a helper function to checkout a DB connection for mutating queries.
func (p *Pool) ForWriting(ctx context.Context, fn func(DB) error) error {
	return pool.With(ctx, p.Pool, func(conn *Conn) error {
		p.wLock.Lock()
		defer p.wLock.Unlock()

		return fn(conn)
	})
}

// Reader is a helper function to checkout a DB connection for read-only queries.
// It returns a DB interface, and a function that must be called to return the
// connection to the pool.
func (p *Pool) Reader(ctx context.Context) (DB, func() error, error) {
	conn, err := p.Pool.Get(ctx)
	if err != nil {
		return nil, nil, err
	}

	p.wLock.RLock()

	close := func() error {
		p.wLock.RUnlock()
		return p.Pool.Put(conn)
	}

	return conn, close, nil
}

var _ DB = (*Pool)(nil)

type poolRowErr struct {
	err error
}

func (r *poolRowErr) Scan(...interface{}) error {
	return r.err
}

func (r *poolRowErr) Err() error {
	return r.err
}

func (r *poolRowErr) Columns() ([]string, error) {
	return nil, r.err
}

var (
	_ Row = (*poolRowErr)(nil)
)
