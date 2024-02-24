package raptor

import (
	"context"

	"github.com/maddiesch/go-raptor/statement/generator"
)

func QueryRowStatement(ctx context.Context, q Querier, stmt generator.Generator) Row {
	query, args, err := stmt.Generate()
	if err != nil {
		return &connRow{err: err}
	}

	return q.QueryRow(ctx, query, args...)
}

func QueryStatement(ctx context.Context, q Querier, stmt generator.Generator) (*Rows, error) {
	query, args, err := stmt.Generate()
	if err != nil {
		return nil, err
	}

	return q.Query(ctx, query, args...)
}

func ExecStatement(ctx context.Context, e Executor, stmt generator.Generator) (Result, error) {
	query, args, err := stmt.Generate()
	if err != nil {
		return nil, err
	}

	return e.Exec(ctx, query, args...)
}

func (c *Conn) QueryRowStatement(ctx context.Context, statement generator.Generator) Row {
	query, args, err := statement.Generate()
	if err != nil {
		return &connRow{err: err}
	}

	return c.QueryRow(ctx, query, args...)
}

func (c *Conn) QueryStatement(ctx context.Context, statement generator.Generator) (*Rows, error) {
	query, args, err := statement.Generate()
	if err != nil {
		return nil, err
	}

	return c.Query(ctx, query, args...)
}

func (c *Conn) ExecStatement(ctx context.Context, statement generator.Generator) (Result, error) {
	query, args, err := statement.Generate()
	if err != nil {
		return nil, err
	}

	return c.Exec(ctx, query, args...)
}
