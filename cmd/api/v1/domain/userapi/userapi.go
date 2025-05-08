package userapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/request"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/logger"
)

type api struct {
	log  *logger.Logger
	core *usercore.Core
}

func (a *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var userReq CreateUserRequest
	if err := request.ParseBody(r, &userReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(userReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	user, err := a.core.Create(ctx, toCreateUserParams(userReq))
	if err != nil {
		if errors.Is(err, usercore.ErrUniqueEmail) {
			return response.WriteError(w, http.StatusConflict, usercore.ErrUniqueEmail)
		}
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	userResp := toUserResponse(user)
	return response.Write(w, http.StatusOK, userResp)
}

func (a *api) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	user, err := getUser(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	userResp := toUserResponse(user)
	return response.Write(w, http.StatusOK, userResp)
}

func (a *api) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	users, err := a.core.Query(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	usersResp := make([]UserResponse, len(users))
	for i, user := range users {
		usersResp[i] = toUserResponse(user)
	}
	return response.Write(w, http.StatusOK, usersResp)
}
