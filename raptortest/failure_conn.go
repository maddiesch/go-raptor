package raptortest

import (
	"context"
	"errors"

	"github.com/maddiesch/go-raptor"
)

// FailureConn is a raptor.DB that always returns an error.
type FailureConn struct{}

func (*FailureConn) Exec(context.Context, string, ...any) (raptor.Result, error) {
	return nil, errors.New("FailureConn.Exec")
}

func (*FailureConn) Query(context.Context, string, ...any) (*raptor.Rows, error) {
	return nil, errors.New("FailureConn.Query")
}

func (*FailureConn) QueryRow(context.Context, string, ...any) raptor.Row {
	return &failureConnRow{}
}

func (f *FailureConn) Transact(_ context.Context, fn func(raptor.DB) error) error {
	return fn(f)
}

type failureConnRow struct{}

func (*failureConnRow) Scan(...any) error {
	return errors.New("FailureConn.Row.Scan")
}

func (*failureConnRow) Err() error {
	return errors.New("FailureConn.Row.Err")
}

func (*failureConnRow) Columns() ([]string, error) {
	return []string{}, errors.New("FailureConn.Row.Columns")
}

var (
	_ raptor.DB  = (*FailureConn)(nil)
	_ raptor.Row = (*failureConnRow)(nil)
)
