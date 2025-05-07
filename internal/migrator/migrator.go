package migrator

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/ardanlabs/darwin/v3"
	"github.com/ardanlabs/darwin/v3/dialects/postgres"
	"github.com/ardanlabs/darwin/v3/drivers/generic"
	"github.com/jmoiron/sqlx"
)

//go:embed sql/migrations.sql
var migrationsDoc string

func Migrate(ctx context.Context, db *sqlx.DB) error {
	driver, err := generic.New(db.DB, postgres.Dialect{})
	if err != nil {
		return fmt.Errorf("construct darwin driver: %w", err)
	}

	d := darwin.New(driver, darwin.ParseMigrations(migrationsDoc))
	return d.Migrate()
}
