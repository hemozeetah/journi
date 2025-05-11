package subscriberdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/subscribercore"
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

func (db *DB) Create(ctx context.Context, subscriber subscribercore.Subscriber) error {
	const q = `
INSERT INTO subscribers
  (user_id, program_id, reference_id, accepted, created_at, updated_at)
VALUES
(:user_id, :program_id, :reference_id, :accepted, :created_at, :updated_at)`

	if err := postgres.ExecContext(ctx, db.db, q, toSubscriberDB(subscriber)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Update(ctx context.Context, subscriber subscribercore.Subscriber) error {
	const q = `
UPDATE subscribers
SET
  accepted = :accepted,
  updated_at = :updated_at
WHERE
  user_id = :user_id AND program_id = :program_id`

	if err := postgres.ExecContext(ctx, db.db, q, toSubscriberDB(subscriber)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, subscriber subscribercore.Subscriber) error {
	const q = `
DELETE FROM subscribers
WHERE
  user_id = :user_id AND program_id = :program_id`

	if err := postgres.ExecContext(ctx, db.db, q, toSubscriberDB(subscriber)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) QueryByID(ctx context.Context, userID uuid.UUID, programID uuid.UUID) (subscribercore.Subscriber, error) {
	const q = `
SELECT
  user_id,
  program_id,
  reference_id,
  accepted,
  created_at,
  updated_at
FROM
  subscribers
WHERE
  user_id = :user_id AND program_id = :program_id`

	data := struct {
		UserID    string `db:"user_id"`
		ProgramID string `db:"program_id"`
	}{
		UserID:    userID.String(),
		ProgramID: programID.String(),
	}

	var s subscriber
	err := postgres.QueryOneContext(ctx, db.db, q, data, &s)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return subscribercore.Subscriber{}, fmt.Errorf("queryonecontext: %w", subscribercore.ErrNotFound)
		}
		return subscribercore.Subscriber{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toSubscriberCore(s), nil
}

func (db *DB) QueryByReferenceID(ctx context.Context, referenceID uuid.UUID) (subscribercore.Subscriber, error) {
	const q = `
SELECT
  user_id,
  program_id,
  reference_id,
  accepted,
  created_at,
  updated_at
FROM
  subscribers
WHERE
  reference_id = :reference_id`

	data := struct {
		ID string `db:"reference_id"`
	}{
		ID: referenceID.String(),
	}

	var s subscriber
	err := postgres.QueryOneContext(ctx, db.db, q, data, &s)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return subscribercore.Subscriber{}, fmt.Errorf("queryonecontext: %w", subscribercore.ErrNotFound)
		}
		return subscribercore.Subscriber{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toSubscriberCore(s), nil
}

func (db *DB) Query(ctx context.Context, query querybuilder.Query) ([]subscribercore.Subscriber, error) {
	const q = `
SELECT
  user_id,
  program_id,
  reference_id,
  accepted,
  created_at,
  updated_at
FROM
  subscribers`

	data := map[string]any{}

	qq := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		q,
		postgres.WhereCluase(fields, query.Constraints),
		postgres.OrderByCluase(fields, query.OrderBy),
		postgres.OffsetCluase(fields, query.Page),
	)

	var ss []subscriber
	err := postgres.QueryContext(ctx, db.db, qq, data, &ss)
	if err != nil {
		return []subscribercore.Subscriber{}, fmt.Errorf("querycontext: %w", err)
	}

	subscribers := make([]subscribercore.Subscriber, len(ss))
	for i, s := range ss {
		subscribers[i] = toSubscriberCore(s)
	}

	return subscribers, nil
}
