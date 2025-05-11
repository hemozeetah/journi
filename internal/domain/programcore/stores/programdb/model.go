package programdb

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/programcore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

var fields = map[querybuilder.Field]string{
	programcore.ID:        "program_id",
	programcore.CompanyID: "company_id",
	programcore.Caption:   "caption",
	programcore.StartDate: "start_date",
	programcore.EndDate:   "end_date",
	programcore.CreatedAt: "created_at",
	programcore.UpdatedAt: "updated_at",
}

type program struct {
	ID        uuid.UUID `db:"program_id"`
	CompanyID uuid.UUID `db:"company_id"`
	Caption   string    `db:"caption"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func toProgramDB(p programcore.Program) program {
	return program{
		ID:        p.ID,
		CompanyID: p.CompanyID,
		Caption:   p.Caption,
		StartDate: p.StartDate,
		EndDate:   p.EndDate,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func toProgramCore(p program) programcore.Program {
	return programcore.Program{
		ID:        p.ID,
		CompanyID: p.CompanyID,
		Caption:   p.Caption,
		StartDate: p.StartDate,
		EndDate:   p.EndDate,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
