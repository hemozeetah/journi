package journeyapi

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/journeycore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

type JourneyResponse struct {
	ID            uuid.UUID `json:"id"`
	ProgramID     uuid.UUID `json:"programID"`
	PlaceID       uuid.UUID `json:"placeID"`
	StartDateTime time.Time `json:"start_datetime"`
	EndDateTime   time.Time `json:"end_datetime"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func toJourneyResponse(journey journeycore.Journey) JourneyResponse {
	return JourneyResponse{
		ID:            journey.ID,
		ProgramID:     journey.ProgramID,
		PlaceID:       journey.PlaceID,
		StartDateTime: journey.StartDateTime,
		EndDateTime:   journey.EndDateTime,
		CreatedAt:     journey.CreatedAt,
		UpdatedAt:     journey.UpdatedAt,
	}
}

type CreateJourneyRequest struct {
	ProgramID     uuid.UUID `json:"programID" validate:"required"`
	PlaceID       uuid.UUID `json:"placeID" validate:"required"`
	StartDateTime time.Time `json:"start_datetime" validate:"required"`
	EndDateTime   time.Time `json:"end_datetime" validate:"required"`
}

func toCreateJourneyParams(journeyReq CreateJourneyRequest) journeycore.CreateJourneyParams {
	return journeycore.CreateJourneyParams{
		ProgramID:     journeyReq.ProgramID,
		PlaceID:       journeyReq.PlaceID,
		StartDateTime: journeyReq.StartDateTime,
		EndDateTime:   journeyReq.EndDateTime,
	}
}

type UpdateJourneyRequest struct {
	PlaceID       *uuid.UUID `json:"placeID" validate:"omitempty,required"`
	StartDateTime *time.Time `json:"start_datetime" validate:"omitempty,required"`
	EndDateTime   *time.Time `json:"end_datetime" validate:"omitempty,required"`
}

func toUpdateJourneyParams(journeyReq UpdateJourneyRequest) journeycore.UpdateJourneyParams {
	return journeycore.UpdateJourneyParams{
		PlaceID:       journeyReq.PlaceID,
		StartDateTime: journeyReq.StartDateTime,
		EndDateTime:   journeyReq.EndDateTime,
	}
}

type params struct {
	ID        string `param:"id" validate:"omitempty,uuid"`
	ProgramID string `param:"program_id" validate:"omitempty,uuid"`
	PlaceID   string `param:"place_id" validate:"omitempty,uuid"`
	OrderBy   string `param:"order_by" validate:"-"`
	Page      string `param:"page" validate:"omitempty,number"`
	Rows      string `param:"rows" validate:"omitempty,number"`
}

var orderByFields = map[string]querybuilder.Field{
	"program_id": journeycore.ProgramID,
	"place_id":   journeycore.PlaceID,
	"start_date": journeycore.StartDateTime,
	"end_date":   journeycore.EndDateTime,
	"created_at": journeycore.CreatedAt,
	"updated_at": journeycore.UpdatedAt,
}

func toQuery(p params) (querybuilder.Query, error) {
	orderBy, err := querybuilder.ParseOrderBy(p.OrderBy, orderByFields, journeycore.DefaultOrderBy)
	if err != nil {
		return querybuilder.Query{}, err
	}

	page, err := querybuilder.ParsePage(p.Page, p.Rows)
	if err != nil {
		return querybuilder.Query{}, err
	}

	constraints := []querybuilder.Constraint{}
	if p.ID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(journeycore.ID, querybuilder.EQ, p.ID))
	}
	if p.ProgramID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(journeycore.ProgramID, querybuilder.EQ, p.ProgramID))
	}
	if p.PlaceID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(journeycore.PlaceID, querybuilder.EQ, p.PlaceID))
	}

	return querybuilder.NewQuery(constraints, orderBy, page), nil
}
