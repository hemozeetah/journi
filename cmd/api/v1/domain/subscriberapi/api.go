package subscriberapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/cmd/api/v1/request"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/programcore"
	"github.com/hemozeetah/journi/internal/domain/subscribercore"
	"github.com/hemozeetah/journi/pkg/logger"
)

type api struct {
	log         *logger.Logger
	core        *subscribercore.Core
	programCore *programcore.Core
}

func (a *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var subscriberReq CreateSubscriberRequest
	if err := request.ParseBody(r, &subscriberReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(subscriberReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	claims, err := jwtauth.GetClaims(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	subscriber, err := a.core.Create(ctx, toCreateSubscriberParams(subscriberReq, claims))
	if err != nil {
		if errors.Is(err, subscribercore.ErrAlreadyExists) {
			return response.WriteError(w, http.StatusConflict, subscribercore.ErrAlreadyExists)
		}
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	subscriberResp := toSubscriberResponse(subscriber)
	return response.Write(w, http.StatusCreated, subscriberResp)
}

func (a *api) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	subscriber, err := getSubscriber(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	subscriberResp := toSubscriberResponse(subscriber)
	return response.Write(w, http.StatusOK, subscriberResp)
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
	subscribers, err := a.core.Query(ctx, query)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	subscribersResp := make([]SubscriberResponse, len(subscribers))
	for i, subscriber := range subscribers {
		subscribersResp[i] = toSubscriberResponse(subscriber)
	}
	return response.Write(w, http.StatusOK, subscribersResp)
}

func (a *api) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var subscriberReq UpdateSubscriberRequest
	if err := request.ParseBody(r, &subscriberReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(subscriberReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	subscriber, err := getSubscriber(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	subscriber, err = a.core.Update(ctx, subscriber, toUpdateSubscriberParams(subscriberReq))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	subscriberResp := toSubscriberResponse(subscriber)
	return response.Write(w, http.StatusOK, subscriberResp)
}

func (a *api) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	subscriber, err := getSubscriber(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	if err := a.core.Delete(ctx, subscriber); err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	return response.Write(w, http.StatusNoContent, nil)
}
