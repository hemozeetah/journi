package api

import (
	"context"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/domain/userapi"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/mux"
	"github.com/hemozeetah/journi/pkg/tracer"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Log *logger.Logger
	DB  *sqlx.DB
}

func New(cfg Config) *mux.Mux {
	mux := mux.New(cfg.Log, generateTraceID(), logging(cfg.Log))

	userapi.Mount(mux, cfg.Log, cfg.DB)

	return mux
}

func generateTraceID() mux.MidFunc {
	return func(handler mux.HandlerFunc) mux.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = tracer.SetRandomID(ctx)
			return handler(ctx, w, r)
		}
	}
}

func logging(log *logger.Logger) mux.MidFunc {
	return func(handler mux.HandlerFunc) mux.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			log.Debug(ctx).
				Attr("method", r.Method).
				Attr("path", r.URL.Path).
				Msg("request")
			return handler(ctx, w, r)
		}
	}
}
