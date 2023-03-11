package test

import (
	"context"

	"github.com/stretchr/testify/require"
)

type TestingT interface {
	require.TestingT

	Log(...any)
	Name() string
}

type TestQueryLogger struct {
	TestingT
}

func (l *TestQueryLogger) LogQuery(_ context.Context, q string, _ []any) {
	l.TestingT.Log("SQL:", q)
}

type CollectedQuery struct {
	Query string
	Args  []any
}

type CollectQueryLogger struct {
	Queries []CollectedQuery
}

func (l *CollectQueryLogger) Reset() {
	l.Queries = nil
}

func (l *CollectQueryLogger) LogQuery(_ context.Context, q string, a []any) {
	l.Queries = append(l.Queries, CollectedQuery{
		Query: q,
		Args:  a,
	})
}
