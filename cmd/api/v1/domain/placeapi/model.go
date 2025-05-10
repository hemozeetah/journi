package placeapi

import (
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/placecore"
)

type PlaceResponse struct {
	ID      uuid.UUID
	CityID  uuid.UUID
	Name    string
	Caption string
	Type    string
	Images  []string
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
	CityID  uuid.UUID
	Name    string
	Caption string
	Type    string
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
	CityID  *uuid.UUID
	Name    *string
	Caption *string
	Type    *string
}

func toUpdatePlaceParams(placeReq UpdatePlaceRequest) placecore.UpdatePlaceParams {
	return placecore.UpdatePlaceParams{
		CityID:  placeReq.CityID,
		Name:    placeReq.Name,
		Caption: placeReq.Caption,
		Type:    placeReq.Type,
	}
}
