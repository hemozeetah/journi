package programcore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/pkg/logger"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

var ErrNotFound = errors.New("program not found")

type Storer interface {
	Create(ctx context.Context, program Program) error
	Update(ctx context.Context, program Program) error
	Delete(ctx context.Context, program Program) error
	QueryByID(ctx context.Context, programID uuid.UUID) (Program, error)
	Query(ctx context.Context, query querybuilder.Query) ([]Program, error)
}

type Core struct {
	log   *logger.Logger
	store Storer
}

func New(log *logger.Logger, store Storer) *Core {
	return &Core{
		log:   log,
		store: store,
	}
}

func (c *Core) Create(ctx context.Context, p CreateProgramParams) (Program, error) {
	now := time.Now()

	program := Program{
		ID:        uuid.New(),
		CompanyID: p.CompanyID,
		StartDate: p.StartDate,
		EndDate:   p.EndDate,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.store.Create(ctx, program); err != nil {
		return Program{}, fmt.Errorf("create: %w", err)
	}

	return program, nil
}

func (c *Core) Update(ctx context.Context, program Program, p UpdateProgramParams) (Program, error) {
	if p.Caption != nil {
		program.Caption = *p.Caption
	}
	if p.StartDate != nil {
		program.StartDate = *p.StartDate
	}
	if p.EndDate != nil {
		program.EndDate = *p.EndDate
	}
	program.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, program); err != nil {
		return Program{}, fmt.Errorf("update: %w", err)
	}

	return program, nil
}

func (c *Core) Delete(ctx context.Context, program Program) error {
	if err := c.store.Delete(ctx, program); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) QueryByID(ctx context.Context, programID uuid.UUID) (Program, error) {
	program, err := c.store.QueryByID(ctx, programID)
	if err != nil {
		return Program{}, fmt.Errorf("querybyid[%s]: %w", programID, err)
	}

	return program, nil
}

func (c *Core) Query(ctx context.Context, query querybuilder.Query) ([]Program, error) {
	programs, err := c.store.Query(ctx, query)
	if err != nil {
		return []Program{}, fmt.Errorf("query: %w", err)
	}

	return programs, nil
}
