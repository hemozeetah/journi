package main

import (
	"context"
	"os"

	"github.com/hemozeetah/journi/pkg/logger"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "12345")

	traceIDFn := func(ctx context.Context) string {
		if traceID, ok := ctx.Value("trace_id").(string); ok {
			return traceID
		}

		return ""
	}
	log := logger.New(os.Stdout, logger.LevelDebug, traceIDFn)

	log.Debug(ctx).
		Attr("foo", "bar").
		Msg("debug message")
}
