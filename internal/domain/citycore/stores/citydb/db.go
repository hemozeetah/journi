package citydb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/citycore"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/postgres"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log *logger.Logger
	db  *sqlx.DB
}

func New(log *logger.Logger, db *sqlx.DB) *DB {
	return &DB{
		log: log,
		db:  db,
	}
}

func (db *DB) Create(ctx context.Context, city citycore.City) error {
	const q = `
INSERT INTO cities
  (city_id, name, caption, images, created_at, updated_at)
VALUES
  (:city_id, :name, :caption, :images, :created_at, :updated_at)`

	if err := postgres.ExecContext(ctx, db.db, q, toCityDB(city)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Update(ctx context.Context, city citycore.City) error {
	const q = `
UPDATE cities
SET
  name = :name,
  caption = :caption,
  images = :images,
  updated_at = :updated_at
WHERE
  city_id = :city_id`

	if err := postgres.ExecContext(ctx, db.db, q, toCityDB(city)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, city citycore.City) error {
	const q = `
DELETE FROM cities
WHERE
  city_id = :city_id`

	if err := postgres.ExecContext(ctx, db.db, q, toCityDB(city)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) QueryByID(ctx context.Context, cityID uuid.UUID) (citycore.City, error) {
	const q = `
SELECT
  city_id,
  name,
  caption,
  images,
  created_at,
  updated_at
FROM
  cities
WHERE
  city_id = :city_id`

	data := struct {
		ID string `db:"city_id"`
	}{
		ID: cityID.String(),
	}

	var c city
	err := postgres.QueryOneContext(ctx, db.db, q, data, &c)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return citycore.City{}, fmt.Errorf("queryonecontext: %w", citycore.ErrNotFound)
		}
		return citycore.City{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toCityCore(c), nil
}

func (db *DB) Query(ctx context.Context) ([]citycore.City, error) {
	const q = `
SELECT
  city_id,
  name,
  caption,
  images,
  created_at,
  updated_at
FROM
  cities`

	data := map[string]any{}

	var cs []city
	err := postgres.QueryContext(ctx, db.db, q, data, &cs)
	if err != nil {
		return []citycore.City{}, fmt.Errorf("querycontext: %w", err)
	}

	cities := make([]citycore.City, len(cs))
	for i, c := range cs {
		cities[i] = toCityCore(c)
	}

	return cities, nil
}
