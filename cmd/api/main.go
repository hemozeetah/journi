package main

import (
	"context"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/mux"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", uuid.UUID{})

	traceIDFn := func(ctx context.Context) string {
		if traceID, ok := ctx.Value("trace_id").(uuid.UUID); ok {
			return traceID.String()
		}

		return ""
	}
	log := logger.New(os.Stdout, logger.LevelDebug, traceIDFn)

	log.Debug(ctx).
		Attr("foo", "bar").
		Msg("debug message")

	app := mux.New(log, generateTraceID)
	app.HandlerFunc("GET", "", "/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		log.Debug(ctx).
			Msg("end point hit")
		w.Write([]byte("hello world"))
		return nil
	})

	http.ListenAndServe(":8080", app)
}

func generateTraceID(handler mux.HandlerFunc) mux.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ctx = context.WithValue(ctx, "trace_id", uuid.New())
		return handler(ctx, w, r)
	}
}
