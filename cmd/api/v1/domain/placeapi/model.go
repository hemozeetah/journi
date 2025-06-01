package placeapi

import (
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/placecore"
	"github.com/hemozeetah/journi/pkg/querybuilder"
)

type PlaceResponse struct {
	ID      uuid.UUID `json:"id"`
	CityID  uuid.UUID `json:"cityID"`
	Name    string    `json:"name"`
	Caption string    `json:"caption"`
	Type    string    `json:"type"`
	Images  []string  `json:"imagesURL"`
}

func toPlaceResponse(place placecore.Place) PlaceResponse {
	images := make([]string, len(place.Images))
	for i, v := range place.Images {
		images[i] = "/static/" + v
	}

	return PlaceResponse{
		ID:      place.ID,
		CityID:  place.CityID,
		Name:    place.Name,
		Caption: place.Caption,
		Type:    place.Type,
		Images:  images,
	}
}

type CreatePlaceRequest struct {
	CityID  uuid.UUID `json:"cityID" validate:"required"`
	Name    string    `json:"name" validate:"required"`
	Caption string    `json:"caption" validate:"required"`
	Type    string    `json:"type" validate:"required"`
}

func toCreatePlaceParams(placeReq CreatePlaceRequest, images []string) placecore.CreatePlaceParams {
	return placecore.CreatePlaceParams{
		CityID:  placeReq.CityID,
		Name:    placeReq.Name,
		Caption: placeReq.Caption,
		Type:    placeReq.Type,
		Images:  images,
	}
}

type UpdatePlaceRequest struct {
	CityID  *uuid.UUID `json:"cityID" validate:"omitempty,required"`
	Name    *string    `json:"name" validate:"omitempty,required"`
	Caption *string    `json:"caption" validate:"omitempty,required"`
	Type    *string    `json:"type" validate:"omitempty,required"`
}

func toUpdatePlaceParams(placeReq UpdatePlaceRequest) placecore.UpdatePlaceParams {
	return placecore.UpdatePlaceParams{
		CityID:  placeReq.CityID,
		Name:    placeReq.Name,
		Caption: placeReq.Caption,
		Type:    placeReq.Type,
	}
}

type params struct {
	ID      string `param:"id" validate:"omitempty,uuid"`
	CityID  string `param:"city_id" validate:"omitempty,uuid"`
	Name    string `param:"name" validate:"omitempty,uuid"`
	Type    string `param:"type" validate:"omitempty,uuid"`
	OrderBy string `param:"order_by" validate:"-"`
	Page    string `param:"page" validate:"omitempty,number"`
	Rows    string `param:"rows" validate:"omitempty,number"`
}

var orderByFields = map[string]querybuilder.Field{
	"created_at": placecore.CreatedAt,
	"updated_at": placecore.UpdatedAt,
}

func toQuery(p params) (querybuilder.Query, error) {
	orderBy, err := querybuilder.ParseOrderBy(p.OrderBy, orderByFields, placecore.DefaultOrderBy)
	if err != nil {
		return querybuilder.Query{}, err
	}

	page, err := querybuilder.ParsePage(p.Page, p.Rows)
	if err != nil {
		return querybuilder.Query{}, err
	}

	constraints := []querybuilder.Constraint{}
	if p.ID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(placecore.ID, querybuilder.EQ, p.ID))
	}
	if p.CityID != "" {
		constraints = append(constraints, querybuilder.NewConstraint(placecore.CityID, querybuilder.EQ, p.CityID))
	}
	if p.Name != "" {
		constraints = append(constraints, querybuilder.NewConstraint(placecore.Name, querybuilder.EQ, p.Name))
	}
	if p.Type != "" {
		constraints = append(constraints, querybuilder.NewConstraint(placecore.Type, querybuilder.EQ, p.Type))
	}

	return querybuilder.NewQuery(constraints, orderBy, page), nil
}
