package placeapi

import (
	"context"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/request"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/placecore"
	"github.com/hemozeetah/journi/pkg/logger"
)

type api struct {
	log  *logger.Logger
	core *placecore.Core
}

func (a *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var placeReq CreatePlaceRequest
	if err := request.ParseBody(r, &placeReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(placeReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	place, err := a.core.Create(ctx, toCreatePlaceParams(placeReq))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	placeResp := toPlaceResponse(place)
	return response.Write(w, http.StatusCreated, placeResp)
}

func (a *api) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	place, err := getPlace(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	placeResp := toPlaceResponse(place)
	return response.Write(w, http.StatusOK, placeResp)
}

func (a *api) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	places, err := a.core.Query(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	placesResp := make([]PlaceResponse, len(places))
	for i, place := range places {
		placesResp[i] = toPlaceResponse(place)
	}
	return response.Write(w, http.StatusOK, placesResp)
}

func (a *api) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var placeReq UpdatePlaceRequest
	if err := request.ParseBody(r, &placeReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(placeReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	place, err := getPlace(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	place, err = a.core.Update(ctx, place, toUpdatePlaceParams(placeReq))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	placeResp := toPlaceResponse(place)
	return response.Write(w, http.StatusOK, placeResp)
}

func (a *api) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	place, err := getPlace(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	if err := a.core.Delete(ctx, place); err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	return response.Write(w, http.StatusNoContent, nil)
}
