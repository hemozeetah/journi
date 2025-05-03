// package tracer provides support for tracing.
package tracer

import (
	"context"

	"github.com/google/uuid"
)

const traceIDKey = "trace_id"

// SetZeroID returns a context with a zero trace id.
func SetZeroID(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceIDKey, uuid.UUID{})
}

// SetRandomID returns a context with a random trace id.
func SetRandomID(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceIDKey, uuid.New())
}

// GetID returns the trace id from given context.
func GetID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(uuid.UUID); ok {
		return traceID.String()
	}

	return ""
}
