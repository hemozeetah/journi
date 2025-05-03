package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/mux"
	"github.com/hemozeetah/journi/pkg/postgres"
	"github.com/spf13/viper"
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

	if err := run(ctx, log); err != nil {
		log.Error(ctx).
			Attr("error", err).
			Msg("failed to run")
		os.Exit(1)
	}
}

type config struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	Database        struct {
		User         string
		Password     string
		Host         string
		Name         string
		DisableTLS   bool
		MaxIdleConns int
		MaxOpenConns int
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	var cfg config

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("readinconfig: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	db, err := postgres.Open(ctx, postgres.Config{
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		Host:         cfg.Database.Host,
		Name:         cfg.Database.Name,
		DisableTLS:   cfg.Database.DisableTLS,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxOpenConns: cfg.Database.MaxOpenConns,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer db.Close()

	log.Info(ctx).
		Attr("status", "database connected").
		Msg("startup")

	app := mux.New(log, generateTraceID)
	app.HandlerFunc("GET", "", "/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		log.Debug(ctx).
			Msg("end point hit")
		w.Write([]byte("hello world"))
		return nil
	})

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      app,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx).
			Attr("status", "server listening").
			Attr("address", server.Addr).
			Msg("startup")

		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx).
			Attr("status", "shutdown started").
			Attr("signal", sig).
			Msg("shutdown")
		defer log.Info(ctx).
			Attr("status", "shutdown completed").
			Attr("signal", sig).
			Msg("shutdown")

		ctx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func generateTraceID(handler mux.HandlerFunc) mux.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ctx = context.WithValue(ctx, "trace_id", uuid.New())
		return handler(ctx, w, r)
	}
}
