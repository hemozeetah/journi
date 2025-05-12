package journeyapi

import (
	"context"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/request"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/journeycore"
	"github.com/hemozeetah/journi/internal/domain/programcore"
	"github.com/hemozeetah/journi/pkg/logger"
)

type api struct {
	log         *logger.Logger
	core        *journeycore.Core
	programCore *programcore.Core
}

func (a *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var journeyReq CreateJourneyRequest
	if err := request.ParseBody(r, &journeyReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(journeyReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	journey, err := a.core.Create(ctx, toCreateJourneyParams(journeyReq))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	journeyResp := toJourneyResponse(journey)
	return response.Write(w, http.StatusCreated, journeyResp)
}

func (a *api) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	journey, err := getJourney(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	journeyResp := toJourneyResponse(journey)
	return response.Write(w, http.StatusOK, journeyResp)
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
	journeys, err := a.core.Query(ctx, query)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	journeysResp := make([]JourneyResponse, len(journeys))
	for i, journey := range journeys {
		journeysResp[i] = toJourneyResponse(journey)
	}
	return response.Write(w, http.StatusOK, journeysResp)
}

func (a *api) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var journeyReq UpdateJourneyRequest
	if err := request.ParseBody(r, &journeyReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(journeyReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	journey, err := getJourney(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	journey, err = a.core.Update(ctx, journey, toUpdateJourneyParams(journeyReq))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	journeyResp := toJourneyResponse(journey)
	return response.Write(w, http.StatusOK, journeyResp)
}

func (a *api) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	journey, err := getJourney(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	if err := a.core.Delete(ctx, journey); err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	return response.Write(w, http.StatusNoContent, nil)
}
