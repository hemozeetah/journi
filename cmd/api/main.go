package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hemozeetah/journi/cmd/api/v1"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/postgres"
	"github.com/hemozeetah/journi/pkg/tracer"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	ctx = tracer.SetZeroID(ctx)

	log := logger.New(os.Stdout, logger.LevelDebug, tracer.GetID)

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

	Database struct {
		User         string
		Password     string
		Host         string
		Name         string
		DisableTLS   bool
		MaxIdleConns int
		MaxOpenConns int
	}

	Auth struct {
		JwtKey    string
		JwtIssuer string
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	var cfg config
	if err := loadEnv(&cfg); err != nil {
		return fmt.Errorf("loading env: %w", err)
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

	app := api.New(api.Config{
		Log:       log,
		DB:        db,
		JwtKey:    cfg.Auth.JwtKey,
		JwtIssuer: cfg.Auth.JwtIssuer,
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

func loadEnv(cfg *config) error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("readinconfig: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}
