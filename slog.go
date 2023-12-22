//go:build go1.21

package raptor

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"
)

type SLogLoggerFunc func(context.Context, string, ...any)

func NewSLogFuncQueryLogger(fn SLogLoggerFunc) QueryLogger {
	return &fnQueryLogger{fn}
}

type fnQueryLogger struct {
	fn SLogLoggerFunc
}

func (l *fnQueryLogger) LogQuery(ctx context.Context, q string, params []any) {
	args := make([]any, len(params))

	for i, a := range params {
		switch a := a.(type) {
		case slog.Attr:
			args[i] = a
		case sql.NamedArg:
			args[i] = slog.Any(a.Name, a.Value)
		case string:
			args[i] = slog.String(strconv.FormatInt(int64(i), 10), a)
		case int64:
			args[i] = slog.Int64(strconv.FormatInt(int64(i), 10), a)
		default:
			args[i] = slog.Any(strconv.FormatInt(int64(i), 10), a)
		}
	}

	l.fn(ctx, q, slog.Group("params", args...))
}
