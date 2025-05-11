package programdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/programcore"
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

func (db *DB) Create(ctx context.Context, program programcore.Program) error {
	const q = `
INSERT INTO programs
  (program_id, company_id, caption, start_date, end_date, created_at, updated_at)
VALUES
  (:program_id, :company_id, :caption, :start_date, :end_date, :created_at, :updated_at)`

	if err := postgres.ExecContext(ctx, db.db, q, toProgramDB(program)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Update(ctx context.Context, program programcore.Program) error {
	const q = `
UPDATE programs
SET
  caption = :caption,
  start_date = :start_date,
  end_date = :end_date,
  updated_at = :updated_at
WHERE
  program_id = :program_id`

	if err := postgres.ExecContext(ctx, db.db, q, toProgramDB(program)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, program programcore.Program) error {
	const q = `
DELETE FROM programs
WHERE
  program_id = :program_id`

	if err := postgres.ExecContext(ctx, db.db, q, toProgramDB(program)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) QueryByID(ctx context.Context, programID uuid.UUID) (programcore.Program, error) {
	const q = `
SELECT
  program_id,
  company_id,
  caption,
  start_date,
  end_date,
  created_at,
  updated_at
FROM
  programs
WHERE
  program_id = :program_id`

	data := struct {
		ID string `db:"program_id"`
	}{
		ID: programID.String(),
	}

	var p program
	err := postgres.QueryOneContext(ctx, db.db, q, data, &p)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return programcore.Program{}, fmt.Errorf("queryonecontext: %w", programcore.ErrNotFound)
		}
		return programcore.Program{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toProgramCore(p), nil
}

func (db *DB) Query(ctx context.Context, query querybuilder.Query) ([]programcore.Program, error) {
	const q = `
SELECT
  program_id,
  company_id,
  caption,
  start_date,
  end_date,
  created_at,
  updated_at
FROM
  programs`

	data := map[string]any{}

	qq := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		q,
		postgres.WhereCluase(fields, query.Constraints),
		postgres.OrderByCluase(fields, query.OrderBy),
		postgres.OffsetCluase(fields, query.Page),
	)

	var ps []program
	err := postgres.QueryContext(ctx, db.db, qq, data, &ps)
	if err != nil {
		return []programcore.Program{}, fmt.Errorf("querycontext: %w", err)
	}

	programs := make([]programcore.Program, len(ps))
	for i, c := range ps {
		programs[i] = toProgramCore(c)
	}

	return programs, nil
}
