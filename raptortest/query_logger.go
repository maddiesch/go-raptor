package raptortest

import (
	"context"

	"github.com/maddiesch/go-raptor"
)

type TLog interface {
	Log(...any)
}

func NewQueryLogger(log TLog) raptor.QueryLogger {
	return &testQueryLogger{log}
}

type testQueryLogger struct {
	TLog
}

func (t *testQueryLogger) LogQuery(_ context.Context, query string, _ []any) {
	t.TLog.Log("Query:", query)
}
