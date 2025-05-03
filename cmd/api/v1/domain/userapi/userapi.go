package userapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/logger"
)

type api struct {
	log  *logger.Logger
	core *usercore.Core
}

func (a *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var userReq CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	user, err := a.core.Create(ctx, toCreateUserParams(userReq))
	if err != nil {
		if errors.Is(err, usercore.ErrUniqueEmail) {
			w.WriteHeader(http.StatusConflict)
			return err
		}

		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	userResp := toUserResponse(user)

	return json.NewEncoder(w).Encode(userResp)
}

func (a *api) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	user, err := getUser(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	userResp := toUserResponse(user)

	return json.NewEncoder(w).Encode(userResp)
}

func (a *api) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	users, err := a.core.Query(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	usersResp := make([]UserResponse, len(users))
	for i, user := range users {
		usersResp[i] = toUserResponse(user)
	}

	return json.NewEncoder(w).Encode(usersResp)
}
