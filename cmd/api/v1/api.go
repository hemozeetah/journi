package api

import (
	"context"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/domain/authapi"
	"github.com/hemozeetah/journi/cmd/api/v1/domain/cityapi"
	"github.com/hemozeetah/journi/cmd/api/v1/domain/journeyapi"
	"github.com/hemozeetah/journi/cmd/api/v1/domain/placeapi"
	"github.com/hemozeetah/journi/cmd/api/v1/domain/postapi"
	"github.com/hemozeetah/journi/cmd/api/v1/domain/programapi"
	"github.com/hemozeetah/journi/cmd/api/v1/domain/userapi"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/muxer"
	"github.com/hemozeetah/journi/pkg/tracer"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Log       *logger.Logger
	DB        *sqlx.DB
	JwtKey    string
	JwtIssuer string
}

func New(cfg Config) *muxer.Mux {
	mux := muxer.New(cfg.Log, tracingMW(), loggingMW(cfg.Log))

	auth := jwtauth.New(jwtauth.Config{
		Log:    cfg.Log,
		DB:     cfg.DB,
		JwtKey: cfg.JwtKey,
		Issuer: cfg.JwtIssuer,
	})

	authapi.Mount(mux, cfg.Log, auth)
	cityapi.Mount(mux, cfg.Log, cfg.DB, auth)
	journeyapi.Mount(mux, cfg.Log, cfg.DB, auth)
	placeapi.Mount(mux, cfg.Log, cfg.DB, auth)
	postapi.Mount(mux, cfg.Log, cfg.DB, auth)
	programapi.Mount(mux, cfg.Log, cfg.DB, auth)
	userapi.Mount(mux, cfg.Log, cfg.DB, auth)

	return mux
}

func tracingMW() muxer.MidFunc {
	return func(handler muxer.HandlerFunc) muxer.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = tracer.SetRandomID(ctx)
			return handler(ctx, w, r)
		}
	}
}

func loggingMW(log *logger.Logger) muxer.MidFunc {
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
