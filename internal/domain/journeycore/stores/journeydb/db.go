package journeydb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/journeycore"
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

func (db *DB) Create(ctx context.Context, journey journeycore.Journey) error {
	const q = `
INSERT INTO journeys
  (journey_id, program_id, place_id, start_datetime, end_datetime, created_at, updated_at)
VALUES
(:journey_id, :program_id, :place_id, :start_datetime, :end_datetime, :created_at, :updated_at)`

	if err := postgres.ExecContext(ctx, db.db, q, toJourneyDB(journey)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Update(ctx context.Context, journey journeycore.Journey) error {
	const q = `
UPDATE journeys
SET
  place_id = :place_id,
  start_datetime = :start_datetime,
  end_datetime = :end_datetime,
  updated_at = :updated_at
WHERE
  journey_id = :journey_id`

	if err := postgres.ExecContext(ctx, db.db, q, toJourneyDB(journey)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, journey journeycore.Journey) error {
	const q = `
DELETE FROM journeys
WHERE
  journey_id = :journey_id`

	if err := postgres.ExecContext(ctx, db.db, q, toJourneyDB(journey)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) QueryByID(ctx context.Context, journeyID uuid.UUID) (journeycore.Journey, error) {
	const q = `
SELECT
  journey_id,
  program_id,
  place_id,
  start_datetime,
  end_datetime,
  created_at,
  updated_at
FROM
  journeys
WHERE
  journey_id = :journey_id`

	data := struct {
		ID string `db:"journey_id"`
	}{
		ID: journeyID.String(),
	}

	var j journey
	err := postgres.QueryOneContext(ctx, db.db, q, data, &j)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return journeycore.Journey{}, fmt.Errorf("queryonecontext: %w", journeycore.ErrNotFound)
		}
		return journeycore.Journey{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toJourneyCore(j), nil
}

func (db *DB) Query(ctx context.Context, query querybuilder.Query) ([]journeycore.Journey, error) {
	const q = `
SELECT
  journey_id,
  program_id,
  place_id,
  start_datetime,
  end_datetime,
  created_at,
  updated_at
FROM
  journeys`

	data := map[string]any{}

	qq := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		q,
		postgres.WhereClause(fields, query.Constraints, data),
		postgres.OrderByClause(fields, query.OrderBy),
		postgres.OffsetClause(fields, query.Page),
	)

	var js []journey
	err := postgres.QueryContext(ctx, db.db, qq, data, &js)
	if err != nil {
		return []journeycore.Journey{}, fmt.Errorf("querycontext: %w", err)
	}

	journeys := make([]journeycore.Journey, len(js))
	for i, j := range js {
		journeys[i] = toJourneyCore(j)
	}

	return journeys, nil
}
