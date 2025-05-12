package programapi

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/programcore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

type ProgramResponse struct {
	ID        uuid.UUID `json:"id"`
	CompanyID uuid.UUID `json:"companyID"`
	Caption   string    `json:"caption"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toProgramResponse(program programcore.Program) ProgramResponse {
	return ProgramResponse{
		ID:        program.ID,
		CompanyID: program.CompanyID,
		Caption:   program.Caption,
		StartDate: program.StartDate,
		EndDate:   program.EndDate,
		CreatedAt: program.CreatedAt,
		UpdatedAt: program.UpdatedAt,
	}
}

type CreateProgramRequest struct {
	Caption   string    `json:"caption" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

func toCreateProgramParams(programReq CreateProgramRequest, claims jwtauth.Claims) programcore.CreateProgramParams {
	return programcore.CreateProgramParams{
		CompanyID: claims.ID,
		Caption:   programReq.Caption,
		StartDate: programReq.StartDate,
		EndDate:   programReq.EndDate,
	}
}

type UpdateProgramRequest struct {
	Caption   *string    `json:"caption" validate:"omitempty,required"`
	StartDate *time.Time `json:"start_date" validate:"omitempty,required"`
	EndDate   *time.Time `json:"end_date" validate:"required"`
}

func toUpdateProgramParams(programReq UpdateProgramRequest) programcore.UpdateProgramParams {
	return programcore.UpdateProgramParams{
		Caption:   programReq.Caption,
		StartDate: programReq.StartDate,
		EndDate:   programReq.EndDate,
	}
}

type params struct {
	ID         string `param:"id" validate:"omitempty,uuid"`
	CompanyID  string `param:"company_id" validate:"omitempty,uuid"`
	StartAfter string `param:"start_after" validate:"omitempty,,datetime=2006-01-02"`
	EndBefore  string `param:"end_before" validate:"omitempty,,datetime=2006-01-02"`
	OrderBy    string `param:"order_by" validate:"-"`
	Page       string `param:"page" validate:"omitempty,number"`
	Rows       string `param:"rows" validate:"omitempty,number"`
}

var orderByFields = map[string]querybuilder.Field{
	"start_date": programcore.StartDate,
	"end_date":   programcore.EndDate,
	"created_at": programcore.CreatedAt,
	"updated_at": programcore.UpdatedAt,
}

func toQuery(p params) (querybuilder.Query, error) {
	orderBy, err := querybuilder.ParseOrderBy(p.OrderBy, orderByFields, programcore.DefaultOrderBy)
	if err != nil {
		return querybuilder.Query{}, err
	}

	page, err := querybuilder.ParsePage(p.Page, p.Rows)
	if err != nil {
		return querybuilder.Query{}, err
	}

	constraints := []querybuilder.Constraint{}
	if p.ID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(programcore.ID, querybuilder.EQ, p.ID))
	}
	if p.CompanyID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(programcore.CompanyID, querybuilder.EQ, p.CompanyID))
	}
	if p.StartAfter != "" {
		constraints = append(constraints, querybuilder.NewConstraint(programcore.StartDate, querybuilder.GTE, p.StartAfter))
	}
	if p.EndBefore != "" {
		constraints = append(constraints, querybuilder.NewConstraint(programcore.EndDate, querybuilder.LTE, p.EndBefore))
	}

	return querybuilder.NewQuery(constraints, orderBy, page), nil
}
