// package logger provides support for initializing the log system.
package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"
)

// TraceIDFn represents a function that return the trace id from given context.
// if the function is nil, the trace id will not be logged.
type TraceIDFn func(ctx context.Context) string

// Level represents a logging level.
type Level slog.Level

// Predefined logging levels.
const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

// Logger represents a logger.
type Logger struct {
	handler   slog.Handler
	traceIDFn TraceIDFn
}

type logger struct {
	Logger
	level slog.Level
	ctx   context.Context
	args  []slog.Attr
}

// New returns a new Logger.
func New(w io.Writer, minLevel Level, traceIDFn TraceIDFn) *Logger {
	// convert time to better format
	// convert source path to only file name with the logged line
	replace := func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey:
			if t, ok := a.Value.Any().(time.Time); ok {
				return slog.Attr{Key: "time", Value: slog.StringValue(t.Format("2006-01-02T15:04:05"))}
			}

		case slog.SourceKey:
			if source, ok := a.Value.Any().(*slog.Source); ok {
				v := fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line)
				return slog.Attr{Key: "file", Value: slog.StringValue(v)}
			}
		}

		return a
	}

	handler := slog.Handler(slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.Level(minLevel),
		ReplaceAttr: replace,
	}))

	return &Logger{
		handler:   handler,
		traceIDFn: traceIDFn,
	}
}

// NewStdLogger returns a standard library Logger that wraps the slog Logger.
func NewStdLogger(logger *Logger, level Level) *log.Logger {
	return slog.NewLogLogger(logger.handler, slog.Level(level))
}

// Info prepares to log at LevelInfo.
func (log *Logger) Info(ctx context.Context) *logger {
	return &logger{
		Logger: *log,
		ctx:    ctx,
		level:  slog.LevelInfo,
	}
}

// Debug prepares to log at LevelDebug.
func (log *Logger) Debug(ctx context.Context) *logger {
	return &logger{
		Logger: *log,
		ctx:    ctx,
		level:  slog.LevelDebug,
	}
}

// Warn prepares to log at LevelWarn.
func (log *Logger) Warn(ctx context.Context) *logger {
	return &logger{
		Logger: *log,
		ctx:    ctx,
		level:  slog.LevelWarn,
	}
}

// Error prepares to log at LevelError.
func (log *Logger) Error(ctx context.Context) *logger {
	return &logger{
		Logger: *log,
		ctx:    ctx,
		level:  slog.LevelError,
	}
}

// Attr adds an attribute for logging.
func (log *logger) Attr(key string, value any) *logger {
	log.args = append(log.args, slog.Any(key, value))
	return log
}

// Msg logs a message with the added attributes.
func (log *logger) Msg(msg string) {
	if !log.handler.Enabled(log.ctx, log.level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])

	r := slog.NewRecord(time.Now(), log.level, msg, pcs[0])

	if log.traceIDFn != nil {
		r.AddAttrs(slog.String("trace_id", log.traceIDFn(log.ctx)))
	}
	r.AddAttrs(log.args...)

	log.handler.Handle(log.ctx, r)
}
