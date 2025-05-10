package placeapi

import (
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/placecore"
)

type PlaceResponse struct {
	ID      uuid.UUID `json:"id"`
	CityID  uuid.UUID `json:"cityID"`
	Name    string    `json:"name"`
	Caption string    `json:"caption"`
	Type    string    `json:"type"`
	Images  []string  `json:"images"`
}

func toPlaceResponse(place placecore.Place) PlaceResponse {
	return PlaceResponse{
		ID:      place.ID,
		CityID:  place.CityID,
		Name:    place.Name,
		Caption: place.Caption,
		Type:    place.Type,
		Images:  place.Images,
	}
}

type CreatePlaceRequest struct {
	CityID  uuid.UUID `json:"cityID" validate:"required"`
	Name    string    `json:"name" validate:"required"`
	Caption string    `json:"caption" validate:"required"`
	Type    string    `json:"type" validate:"required"`
}

func toCreatePlaceParams(placeReq CreatePlaceRequest) placecore.CreatePlaceParams {
	return placecore.CreatePlaceParams{
		CityID:  placeReq.CityID,
		Name:    placeReq.Name,
		Caption: placeReq.Caption,
		Type:    placeReq.Type,
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
