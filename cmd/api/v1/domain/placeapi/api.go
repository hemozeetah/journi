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
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	var placeReq CreatePlaceRequest
	if err := request.ParseForm(r, &placeReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(placeReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	images, err := request.ParseFile(r, "images")
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	place, err := a.core.Create(ctx, toCreatePlaceParams(placeReq, images))
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
	var p params
	if err := request.ParseQueryParams(r, &p); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(p); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	query, err := toQuery(p)
	if err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	places, err := a.core.Query(ctx, query)
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
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	var placeReq UpdatePlaceRequest
	if err := request.ParseForm(r, &placeReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(placeReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	images, err := request.ParseFile(r, "images")
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	place, err := getPlace(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	place, err = a.core.Update(ctx, place, toUpdatePlaceParams(placeReq, images))
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
