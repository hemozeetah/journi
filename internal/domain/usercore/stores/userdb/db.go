package userdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/usercore"
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

func (db *DB) Create(ctx context.Context, user usercore.User) error {
	const q = `
INSERT INTO users
  (user_id, name, email, password, role, profile, created_at, updated_at)
VALUES
  (:user_id, :name, :email, :password, :role, :profile, :created_at, :updated_at)`

	if err := postgres.ExecContext(ctx, db.db, q, toUserDB(user)); err != nil {
		if errors.Is(err, postgres.ErrDBDuplicatedEntry) {
			return fmt.Errorf("execcontext: %w", usercore.ErrUniqueEmail)
		}
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Update(ctx context.Context, user usercore.User) error {
	const q = `
UPDATE users
SET
  name = :name,
  email = :email,
  password = :password,
  role = :role,
  profile = :profile,
  updated_at = :updated_at
WHERE
  user_id = :user_id`

	if err := postgres.ExecContext(ctx, db.db, q, toUserDB(user)); err != nil {
		if errors.Is(err, postgres.ErrDBDuplicatedEntry) {
			return fmt.Errorf("execcontext: %w", usercore.ErrUniqueEmail)
		}
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, user usercore.User) error {
	const q = `
DELETE FROM users
WHERE
  user_id = :user_id`

	if err := postgres.ExecContext(ctx, db.db, q, toUserDB(user)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (db *DB) QueryByID(ctx context.Context, userID uuid.UUID) (usercore.User, error) {
	const q = `
SELECT
  user_id,
  name,
  email,
  password,
  role,
  profile,
  created_at,
  updated_at
FROM
  users
WHERE
  user_id = :user_id`

	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	var u user
	err := postgres.QueryOneContext(ctx, db.db, q, data, &u)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return usercore.User{}, fmt.Errorf("queryonecontext: %w", usercore.ErrNotFound)
		}
		return usercore.User{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toUserCore(u), nil
}

func (db *DB) QueryByEmail(ctx context.Context, email string) (usercore.User, error) {
	const q = `
SELECT
  user_id,
  name,
  email,
  password,
  role,
  profile,
  created_at,
  updated_at
FROM
  users
WHERE
  email = :email`

	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	var u user
	err := postgres.QueryOneContext(ctx, db.db, q, data, &u)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return usercore.User{}, fmt.Errorf("queryonecontext: %w", usercore.ErrNotFound)
		}
		return usercore.User{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toUserCore(u), nil
}

func (db *DB) Query(ctx context.Context) ([]usercore.User, error) {
	const q = `
SELECT
  user_id,
  name,
  email,
  password,
  role,
  profile,
  created_at,
  updated_at
FROM
  users`

	data := map[string]any{}

	var us []user
	err := postgres.QueryContext(ctx, db.db, q, data, &us)
	if err != nil {
		return []usercore.User{}, fmt.Errorf("querycontext: %w", err)
	}

	users := make([]usercore.User, len(us))
	for i, u := range us {
		users[i] = toUserCore(u)
	}

	return users, nil
}
