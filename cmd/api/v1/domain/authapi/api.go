package authapi

import (
	"context"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/auth"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/pkg/logger"
)

type api struct {
	log  *logger.Logger
	auth *auth.Auth
}

func (a api) token(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	token, err := a.auth.GenerateToken(claims)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	resp := response.Envelope{
		"token":  token,
		"claims": claims,
	}
	return response.Write(w, http.StatusOK, resp)
}

func (a api) claims(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	resp := response.Envelope{
		"claims": claims,
	}
	return response.Write(w, http.StatusOK, resp)
}

func (a api) authorizeAdmin(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	resp := response.Envelope{
		"claims":  claims,
		"isAdmin": true,
	}
	return response.Write(w, http.StatusOK, resp)
}
