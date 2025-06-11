package postapi

import (
	"time"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/postcore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

type PostResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userID"`
	PlaceID   uuid.UUID `json:"placeID"`
	Caption   string    `json:"caption"`
	Images    []string  `json:"imagesURL"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toPostResponse(post postcore.Post) PostResponse {
	images := make([]string, len(post.Images))
	for i, v := range post.Images {
		images[i] = "/static/" + v
	}

	return PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		PlaceID:   post.PlaceID,
		Caption:   post.Caption,
		Images:    images,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

type CreatePostRequest struct {
	PlaceID uuid.UUID `json:"placeID" validate:"required"`
	Caption string    `json:"caption" validate:"required"`
}

func toCreatePostParams(postReq CreatePostRequest, claims jwtauth.Claims, images []string) postcore.CreatePostParams {
	return postcore.CreatePostParams{
		UserID:  claims.ID,
		PlaceID: postReq.PlaceID,
		Caption: postReq.Caption,
		Images:  images,
	}
}

type UpdatePostRequest struct {
	PlaceID *uuid.UUID `json:"placeID" validate:"omitempty,required"`
	Caption *string    `json:"caption" validate:"omitempty,required"`
}

func toUpdatePostParams(postReq UpdatePostRequest) postcore.UpdatePostParams {
	return postcore.UpdatePostParams{
		PlaceID: postReq.PlaceID,
		Caption: postReq.Caption,
	}
}

type params struct {
	ID      string `param:"id" validate:"omitempty,uuid"`
	UserID  string `param:"user_id" validate:"omitempty,uuid"`
	PlaceID string `param:"place_id" validate:"omitempty,uuid"`
	OrderBy string `param:"order_by" validate:"-"`
	Page    string `param:"page" validate:"omitempty,number"`
	Rows    string `param:"rows" validate:"omitempty,number"`
}

var orderByFields = map[string]querybuilder.Field{
	"created_at": postcore.CreatedAt,
	"updated_at": postcore.UpdatedAt,
}

func toQuery(p params) (querybuilder.Query, error) {
	orderBy, err := querybuilder.ParseOrderBy(p.OrderBy, orderByFields, postcore.DefaultOrderBy)
	if err != nil {
		return querybuilder.Query{}, err
	}

	page, err := querybuilder.ParsePage(p.Page, p.Rows)
	if err != nil {
		return querybuilder.Query{}, err
	}

	constraints := []querybuilder.Constraint{}
	if p.ID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(postcore.ID, querybuilder.EQ, p.ID))
	}
	if p.UserID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(postcore.UserID, querybuilder.EQ, p.UserID))
	}
	if p.PlaceID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(postcore.PlaceID, querybuilder.EQ, p.PlaceID))
	}

	return querybuilder.NewQuery(constraints, orderBy, page), nil
}
