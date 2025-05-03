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

type Store struct {
	log *logger.Logger
	db  *sqlx.DB
}

func New(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

func (s *Store) Create(ctx context.Context, user usercore.User) error {
	const q = `
INSERT INTO users
  (id, name, email, password)
VALUES
  (:id, :name, :email, :password)`

	if err := postgres.ExecContext(ctx, s.db, q, toUserDB(user)); err != nil {
		if errors.Is(err, postgres.ErrDBDuplicatedEntry) {
			return fmt.Errorf("execcontext: %w", usercore.ErrUniqueEmail)
		}
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, user usercore.User) error {
	const q = `
UPDATE users
SET
  name = :name,
  email = :email,
  password = :password
WHERE
  id = :id`

	if err := postgres.ExecContext(ctx, s.db, q, toUserDB(user)); err != nil {
		if errors.Is(err, postgres.ErrDBDuplicatedEntry) {
			return fmt.Errorf("execcontext: %w", usercore.ErrUniqueEmail)
		}
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, user usercore.User) error {
	const q = `
DELETE FROM users
WHERE
  id = :id`

	if err := postgres.ExecContext(ctx, s.db, q, toUserDB(user)); err != nil {
		return fmt.Errorf("execcontext: %w", err)
	}

	return nil
}

func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (usercore.User, error) {
	const q = `
SELECT
  id,
  name,
  email,
  password
FROM
  users
WHERE
  id = :id`

	data := struct {
		ID string `db:"id"`
	}{
		ID: userID.String(),
	}

	var u user
	err := postgres.QueryOneContext(ctx, s.db, q, data, &u)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return usercore.User{}, fmt.Errorf("queryonecontext: %w", usercore.ErrNotFound)
		}
		return usercore.User{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toUserCore(u), nil
}

func (s *Store) QueryByEmail(ctx context.Context, email string) (usercore.User, error) {
	const q = `
SELECT
  id,
  name,
  email,
  password
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
	err := postgres.QueryOneContext(ctx, s.db, q, data, &u)
	if err != nil {
		if errors.Is(err, postgres.ErrDBNotFound) {
			return usercore.User{}, fmt.Errorf("queryonecontext: %w", usercore.ErrNotFound)
		}
		return usercore.User{}, fmt.Errorf("queryonecontext: %w", err)
	}

	return toUserCore(u), nil
}

func (s *Store) Query(ctx context.Context) ([]usercore.User, error) {
	const q = `
SELECT
  id,
  name,
  email,
  password
FROM
  users`

	data := map[string]any{}

	var us []user
	err := postgres.QueryContext(ctx, s.db, q, data, &us)
	if err != nil {
		return []usercore.User{}, fmt.Errorf("querycontext: %w", err)
	}

	users := make([]usercore.User, len(us))
	for i, u := range us {
		users[i] = toUserCore(u)
	}

	return users, nil
}
