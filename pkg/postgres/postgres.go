// package postgres provides support for connecting to postgres databases.
package postgres

import (
	"context"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config represents a configuration for connecting to a postgres database.
type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	DisableTLS   bool
	MaxIdleConns int
	MaxOpenConns int
}

// Open returns a connection to the postgres database.
func Open(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	q := make(url.Values)
	if cfg.DisableTLS {
		q.Set("sslmode", "disable")
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
