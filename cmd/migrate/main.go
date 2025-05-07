package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hemozeetah/journi/internal/migrator"
	"github.com/hemozeetah/journi/pkg/postgres"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
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

func run(ctx context.Context) error {
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

	if err := migrator.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrating db: %w", err)
	}

	fmt.Println("database migrated successfully")

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
