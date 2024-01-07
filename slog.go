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
		args[i] = CreateSlogArg(int64(i), a)
	}

	l.fn(ctx, q, slog.Group("params", args...))
}

func CreateSlogArg(i int64, v any) slog.Attr {
	switch v := v.(type) {
	case slog.Attr:
		return v
	case sql.NamedArg:
		return slog.Any(v.Name, v.Value)
	case string:
		return slog.String(strconv.FormatInt(int64(i), 10), v)
	case int:
		return slog.Int(strconv.FormatInt(int64(i), 10), v)
	case int64:
		return slog.Int64(strconv.FormatInt(int64(i), 10), v)
	case uint:
		return slog.Uint64(strconv.FormatInt(int64(i), 10), uint64(v))
	case uint64:
		return slog.Uint64(strconv.FormatInt(int64(i), 10), v)
	default:
		return slog.Any(strconv.FormatInt(int64(i), 10), v)
	}
}
