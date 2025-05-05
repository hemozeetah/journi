package api

import (
	"context"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/domain/userapi"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/hemozeetah/journi/pkg/tracer"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Log *logger.Logger
	DB  *sqlx.DB
}

func New(cfg Config) *muxer.Mux {
	mux := muxer.New(cfg.Log, generateTraceID(), logging(cfg.Log))

	userapi.Mount(mux, cfg.Log, cfg.DB)

	return mux
}

func generateTraceID() muxer.MidFunc {
	return func(handler muxer.HandlerFunc) muxer.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = tracer.SetRandomID(ctx)
			return handler(ctx, w, r)
		}
	}
}

func logging(log *logger.Logger) muxer.MidFunc {
	return func(handler muxer.HandlerFunc) muxer.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			log.Debug(ctx).
				Attr("method", r.Method).
				Attr("path", r.URL.Path).
				Msg("request")
			return handler(ctx, w, r)
		}
	}
}
