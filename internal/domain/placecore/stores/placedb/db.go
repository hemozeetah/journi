package placedb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/placecore"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/postgres"
	"github.com/hemozeetah/journi/pkg/querybuilder"
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

func (db *DB) Create(ctx context.Context, place placecore.Place) error {
	const q = `
INSERT INTO places
  (place_id, city_id, name, caption, type, images, created_at, updated_at)
VALUES
  (:place_id, :city_id, :name, :caption, :type, :images, :created_at, :updated_at)`

	if err := postgres.ExecContext(ctx, db.db, q, toPlaceDB(place)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Update(ctx context.Context, place placecore.Place) error {
	const q = `
UPDATE places
SET
  city_id = :city_id,
  name = :name,
  caption = :caption,
  type = :type,
  images = :images,
  updated_at = :updated_at
WHERE
  place_id = :place_id`

	if err := postgres.ExecContext(ctx, db.db, q, toPlaceDB(place)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, place placecore.Place) error {
	const q = `
DELETE FROM places
WHERE
  place_id = :place_id`

	if err := postgres.ExecContext(ctx, db.db, q, toPlaceDB(place)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) QueryByID(ctx context.Context, placeID uuid.UUID) (placecore.Place, error) {
	const q = `
SELECT
  place_id,
  city_id,
  name,
  caption,
  type,
  images,
  created_at,
  updated_at
FROM
  places
WHERE
  place_id = :place_id`

	data := struct {
		ID string `db:"place_id"`
	}{
		ID: placeID.String(),
	}

	var p place
	err := postgres.QueryOneContext(ctx, db.db, q, data, &p)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return placecore.Place{}, fmt.Errorf("queryonecontext: %w", placecore.ErrNotFound)
		}
		return placecore.Place{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toPlaceCore(p), nil
}

func (db *DB) Query(ctx context.Context, query querybuilder.Query) ([]placecore.Place, error) {
	const q = `
SELECT
  place_id,
  city_id,
  name,
  caption,
  type,
  images,
  created_at,
  updated_at
FROM
  places`

	data := map[string]any{}

	qq := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		q,
		postgres.WhereClause(fields, query.Constraints, data),
		postgres.OrderByClause(fields, query.OrderBy),
		postgres.OffsetClause(fields, query.Page),
	)

	var ps []place
	err := postgres.QueryContext(ctx, db.db, qq, data, &ps)
	if err != nil {
		return []placecore.Place{}, fmt.Errorf("querycontext: %w", err)
	}

	places := make([]placecore.Place, len(ps))
	for i, p := range ps {
		places[i] = toPlaceCore(p)
	}

	return places, nil
}
