package raptor

import (
	"context"

	"github.com/maddiesch/go-raptor/statement/generator"
)

func (c *Conn) QueryRowStatement(ctx context.Context, statement generator.Generator) *Row {
	query, args, err := statement.Generate()
	if err != nil {
		return &Row{err: err}
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
