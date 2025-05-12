package subscriberapi

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/subscribercore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

type SubscriberResponse struct {
	UserID      uuid.UUID `json:"id"`
	ProgramID   uuid.UUID `json:"programID"`
	ReferenceID uuid.UUID `json:"referenceID"`
	Accepted    bool      `json:"accepted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func toSubscriberResponse(subscriber subscribercore.Subscriber) SubscriberResponse {
	return SubscriberResponse{
		UserID:      subscriber.UserID,
		ProgramID:   subscriber.ProgramID,
		ReferenceID: subscriber.ReferenceID,
		CreatedAt:   subscriber.CreatedAt,
		UpdatedAt:   subscriber.UpdatedAt,
	}
}

type CreateSubscriberRequest struct {
	ProgramID uuid.UUID `json:"programID" validate:"required"`
}

func toCreateSubscriberParams(subscriberReq CreateSubscriberRequest, claims jwtauth.Claims) subscribercore.CreateSubscriberParams {
	return subscribercore.CreateSubscriberParams{
		UserID:    claims.ID,
		ProgramID: subscriberReq.ProgramID,
	}
}

type UpdateSubscriberRequest struct {
	Accepted bool `json:"accepted"`
}

func toUpdateSubscriberParams(subscriberReq UpdateSubscriberRequest) subscribercore.UpdateSubscriberParams {
	return subscribercore.UpdateSubscriberParams{
		Accepted: subscriberReq.Accepted,
	}
}

type params struct {
	UserID    string `param:"user_id" validate:"omitempty,uuid"`
	ProgramID string `param:"program_id" validate:"omitempty,uuid"`
	Accepted  string `param:"accepted" validate:"omitempty,boolean"`
	OrderBy   string `param:"order_by" validate:"-"`
	Page      string `param:"page" validate:"omitempty,number"`
	Rows      string `param:"rows" validate:"omitempty,number"`
}

var orderByFields = map[string]querybuilder.Field{
	"accepted":   subscribercore.Accepted,
	"created_at": subscribercore.CreatedAt,
	"updated_at": subscribercore.UpdatedAt,
}

func toQuery(p params) (querybuilder.Query, error) {
	orderBy, err := querybuilder.ParseOrderBy(p.OrderBy, orderByFields, subscribercore.DefaultOrderBy)
	if err != nil {
		return querybuilder.Query{}, err
	}

	page, err := querybuilder.ParsePage(p.Page, p.Rows)
	if err != nil {
		return querybuilder.Query{}, err
	}

	constraints := []querybuilder.Constraint{}
	if p.UserID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(subscribercore.UserID, querybuilder.EQ, p.UserID))
	}
	if p.ProgramID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(subscribercore.ProgramID, querybuilder.EQ, p.ProgramID))
	}
	if p.Accepted != "" {
		constraints = append(constraints, querybuilder.NewConstraint(subscribercore.Accepted, querybuilder.GTE, p.Accepted))
	}

	return querybuilder.NewQuery(constraints, orderBy, page), nil
}
