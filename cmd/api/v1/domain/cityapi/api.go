package cityapi

import (
	"context"
	"embed"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/request"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/citycore"
	"github.com/hemozeetah/journi/pkg/logger"
)

type api struct {
	log    *logger.Logger
	core   *citycore.Core
	static *embed.FS
}

func (a *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	var cityReq CreateCityRequest
	if err := request.ParseForm(r, &cityReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(cityReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	images, err := request.ParseFile(r, "images")
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	city, err := a.core.Create(ctx, toCreateCityParams(cityReq, images))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	cityResp := toCityResponse(city)
	return response.Write(w, http.StatusCreated, cityResp)
}

func (a *api) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	city, err := getCity(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	cityResp := toCityResponse(city)
	return response.Write(w, http.StatusOK, cityResp)
}

func (a *api) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cities, err := a.core.Query(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	citiesResp := make([]CityResponse, len(cities))
	for i, city := range cities {
		citiesResp[i] = toCityResponse(city)
	}
	return response.Write(w, http.StatusOK, citiesResp)
}

func (a *api) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	var cityReq UpdateCityRequest
	if err := request.ParseForm(r, &cityReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(cityReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	images, err := request.ParseFile(r, "images")
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	city, err := getCity(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	city, err = a.core.Update(ctx, city, toUpdateCityParams(cityReq, images))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	cityResp := toCityResponse(city)
	return response.Write(w, http.StatusOK, cityResp)
}

func (a *api) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	city, err := getCity(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	if err := a.core.Delete(ctx, city); err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	return response.Write(w, http.StatusNoContent, nil)
}
