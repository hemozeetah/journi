package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// lib/pq errorCodeNames
// https://github.com/lib/pq/blob/master/error.go
const (
	uniqueViolation = "23505"
	undefinedTable  = "42P01"
)

// Set of error variables for CRUD operations.
var (
	ErrDBNotFound        = sql.ErrNoRows
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
	ErrUndefinedTable    = errors.New("undefined table")
)

// ExecContext executes a query.
func ExecContext(ctx context.Context, db *sqlx.DB, query string, arg any) error {
	if _, err := db.NamedExecContext(ctx, query, arg); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case uniqueViolation:
				return ErrDBDuplicatedEntry

			case undefinedTable:
				return ErrUndefinedTable
			}
		}

		return err
	}

	return nil
}

// QueryOneContext executes a query.
func QueryOneContext(ctx context.Context, db *sqlx.DB, query string, arg any, dest any) error {
	rows, err := db.NamedQueryContext(ctx, query, arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == undefinedTable {
			return ErrUndefinedTable
		}

		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrDBNotFound
	}

	err = rows.StructScan(dest)
	if err != nil {
		return err
	}

	return nil
}

// QueryContext executes a query.
func QueryContext[T any](ctx context.Context, db *sqlx.DB, query string, arg any, dest *[]T) error {
	rows, err := db.NamedQueryContext(ctx, query, arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == undefinedTable {
			return ErrUndefinedTable
		}

		return err
	}
	defer rows.Close()

	var slice []T
	for rows.Next() {
		v := new(T)
		if err = rows.StructScan(v); err != nil {
			return err
		}

		slice = append(slice, *v)
	}

	*dest = slice

	return nil
}
