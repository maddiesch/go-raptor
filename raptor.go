// Package raptor provides a simple and easy-to-use interface for working with SQLite3 databases.
package raptor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	_ "modernc.org/sqlite"
)

const (
	// DriverName is the name of the SQLite3 driver.
	DriverName = "sqlite"
)

var (
	connID atomic.Uint64
)

// New opens a new database connection
func New(source string) (*Conn, error) {
	db, err := sql.Open(DriverName, source)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		db:       db,
		id:       connID.Add(1),
		queryLog: new(atomic.Pointer[logger]),
	}
	c.SetLogger(NewNoopQueryLogger())

	return c, nil
}

// Conn represents a connection to a SQLite3 database.
type Conn struct {
	id       uint64                  // Connection id
	sp       atomic.Uint64           // Savepoint id
	db       *sql.DB                 // Underlying database connection
	queryLog *atomic.Pointer[logger] // Log query
}

// Close the database connection and perform any necessary cleanup
//
// Once close is called, new queries will be rejected.
// Close will block until all outstanding queries have completed.
func (c *Conn) Close() error {
	return c.db.Close()
}

type logger struct {
	QueryLogger
}

func (c *Conn) SetLogger(l QueryLogger) {
	c.queryLog.Store(&logger{l})
}

// SetLogger assigns a logger instance to the connection.
func (c *Conn) SetQueryLogger(l QueryLogger) {
	c.queryLog.Store(&logger{l})
}

func (c *Conn) queryLogger() QueryLogger {
	return c.queryLog.Load().QueryLogger
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
func (c *Conn) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// QueryLogger provides a standard interface for logging all SQL queries sent to Raptor
type QueryLogger interface {
	LogQuery(context.Context, string, []any)
}

// NewQueryLogger creates a new QueryLogger that logs queries to an io.Writer.
func NewQueryLogger(w io.Writer) QueryLogger {
	return &wQueryLogger{w}
}

type wQueryLogger struct {
	w io.Writer
}

func (w *wQueryLogger) LogQuery(_ context.Context, q string, _ []any) {
	fmt.Fprintln(w.w, q)
}

type noopQueryLogger struct{}

// NewNoopQueryLogger creates a new QueryLogger that doesn't log any queries.
func NewNoopQueryLogger() QueryLogger {
	return &noopQueryLogger{}
}

func (w *noopQueryLogger) LogQuery(context.Context, string, []any) {}

// A Result summarizes an executed SQL command.
type Result interface {
	sql.Result
}

// Rows is the result of a query. See sql.Rows for more information.
type Rows struct {
	*sql.Rows
}

var (
	ErrNoRows = sql.ErrNoRows
)

// Row is the result of calling QueryRow to select a single row.
type Row interface {
	Scanner

	Err() error
}

type connRow struct {
	rows *sql.Rows
	err  error
}

func (r *connRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	defer r.rows.Close()
	for _, dp := range dest {
		if _, ok := dp.(*sql.RawBytes); ok {
			return errors.New("raptor: RawBytes isn't allowed on Row.Scan")
		}
	}

	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return err
		}
		return ErrNoRows
	}

	err := r.rows.Scan(dest...)
	if err != nil {
		return err
	}
	// Make sure the query can be processed to completion with no errors.
	return r.rows.Close()
}

func (r *connRow) Err() error {
	return r.err
}

func (r *connRow) Columns() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.rows.Columns()
}

type Scanner interface {
	Scan(...any) error
	Columns() ([]string, error)
}

// Executor defines an interface for executing queries that don't return rows.
type Executor interface {
	Exec(context.Context, string, ...any) (Result, error)
}

// Exec perform a query on the database. It will not return any rows. e.g. insert or delete
func (c *Conn) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	return c.exec(ctx, query, args...)
}

func (c *Conn) exec(ctx context.Context, query string, args ...any) (Result, error) {
	c.queryLogger().LogQuery(ctx, query, args)

	r, err := c.db.ExecContext(ctx, query, args...)

	return Result(r), err
}

// Querier defines an interface for executing queries that return rows from the database.
type Querier interface {
	Query(context.Context, string, ...any) (*Rows, error)
	QueryRow(context.Context, string, ...any) Row
}

func (c *Conn) Query(ctx context.Context, query string, args ...any) (*Rows, error) {
	return c.query(ctx, query, args)
}

func (c *Conn) query(ctx context.Context, query string, args []any) (*Rows, error) {
	c.queryLogger().LogQuery(ctx, query, args)

	r, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &Rows{r}, nil
}

func (c *Conn) QueryRow(ctx context.Context, query string, args ...any) Row {
	return c.queryRow(ctx, query, args)
}

func (c *Conn) queryRow(ctx context.Context, query string, args []any) Row {
	c.queryLogger().LogQuery(ctx, query, args)

	r, err := c.db.QueryContext(ctx, query, args...)

	return &connRow{rows: r, err: err}
}

func (c *Conn) newSavepointName() string {
	return fmt.Sprintf("tx_%d_%d", c.id, c.sp.Add(1))
}

// TxRollbackError is returned when a transaction is rolled back and the rollback also returns an error.
type TxRollbackError struct {
	Underlying error
	Rollback   error
}

func (e *TxRollbackError) Error() string {
	return fmt.Sprintf("rollback error: %s; rollback error: %s", e.Underlying, e.Rollback)
}

// TxBroker defines an interface for performing a transaction.
type TxBroker interface {
	Transact(context.Context, func(DB) error) error
}

// DB defines a standard set of interfaces that allow CRUD operations on a database.
type DB interface {
	Executor
	Querier
	TxBroker
}

var _ DB = (*Conn)(nil)
var _ DB = (*txConn)(nil)

func (c *Conn) Transact(ctx context.Context, fn func(DB) error) error {
	return c.transact(ctx, 0, fn)
}

func (c *Conn) transact(ctx context.Context, depth int, fn func(DB) error) error {
	savepoint := c.newSavepointName()

	txConn := &txConn{
		conn:  c,
		depth: depth,
		name:  savepoint,
		state: txStateInit,
	}

	if err := txConn.begin(ctx); err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			txConn.rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(txConn); err != nil {
		if rErr := txConn.rollback(ctx); rErr != nil {
			return &TxRollbackError{Underlying: err, Rollback: rErr}
		}
		if errors.Is(err, ErrTxRollback) {
			return nil
		}
		return err
	}

	return txConn.commit(ctx)
}

const (
	txStateInit uint8 = iota
	txStateRunning
	txStateCommitted
	txStateRollbacked
)

type txConn struct {
	conn  *Conn
	depth int
	mu    sync.Mutex
	name  string
	state uint8
}

var (
	ErrTransactionAlreadyStarted = errors.New("transaction already started")
	ErrTransactionNotRunning     = errors.New("transaction not running")
	ErrTxRollback                = errors.New("transaction rollback") // Can be returned from a transaction to rollback the transaction. Will not be returned to the caller
)

func (t *txConn) begin(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != txStateInit {
		return ErrTransactionAlreadyStarted
	}

	_, err := t.conn.exec(ctx, "SAVEPOINT "+t.name+";")
	if err == nil {
		t.state = 1
	}

	return err
}

func (t *txConn) rollback(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != txStateRunning {
		return nil
	}

	_, err := t.conn.exec(ctx, "ROLLBACK TRANSACTION TO SAVEPOINT "+t.name+";")
	if err == nil {
		t.state = txStateRollbacked
	}

	return nil
}

func (t *txConn) commit(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != txStateRunning {
		return nil
	}

	_, err := t.conn.exec(ctx, "RELEASE SAVEPOINT "+t.name+";")
	if err == nil {
		t.state = txStateCommitted
	}

	return err
}

func (t *txConn) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != txStateRunning {
		return nil, ErrTransactionNotRunning
	}

	return t.conn.exec(ctx, query, args...)
}

func (t *txConn) Transact(ctx context.Context, fn func(DB) error) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != txStateRunning {
		return ErrTransactionNotRunning
	}

	return t.conn.transact(ctx, t.depth+1, fn)
}

func (t *txConn) Query(ctx context.Context, query string, args ...any) (*Rows, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != txStateRunning {
		return nil, ErrTransactionNotRunning
	}

	return t.conn.query(ctx, query, args)
}

func (t *txConn) QueryRow(ctx context.Context, query string, args ...any) Row {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != txStateRunning {
		return &connRow{rows: nil, err: ErrTransactionNotRunning}
	}

	return t.conn.queryRow(ctx, query, args)
}
