package jwtauth

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/muxer"
)

const claimsKey = "claims"

func setClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func GetClaims(ctx context.Context) (Claims, error) {
	v, ok := ctx.Value(claimsKey).(Claims)
	if !ok {
		return Claims{}, errors.New("user id not found in context")
	}

	return v, nil
}

func (a *Auth) BasicMW() muxer.MidFunc {
	return func(next muxer.HandlerFunc) muxer.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			email, pass, err := parseBasicAuth(r.Header.Get("authorization"))
			if err != nil {
				return response.WriteError(w, http.StatusUnauthorized, err)
			}

			user, err := a.core.Authenticate(ctx, email, pass)
			if err != nil {
				switch {
				case errors.Is(err, usercore.ErrAuthenticationFailure):
					return response.WriteError(w, http.StatusUnauthorized, usercore.ErrAuthenticationFailure)

				case errors.Is(err, usercore.ErrNotFound):
					return response.WriteError(w, http.StatusNotFound, usercore.ErrNotFound)

				default:
					return response.WriteError(w, http.StatusInternalServerError, err)
				}
			}

			claims := a.newClaims(user)
			ctx = setClaims(ctx, claims)

			return next(ctx, w, r)
		}
	}
}

func parseBasicAuth(auth string) (string, string, error) {
	if !strings.HasPrefix(auth, "Basic ") {
		return "", "", errors.New("expected authorization header format: Basic <auth>")
	}

	auth = auth[6:]

	c, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return "", "", err
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", errors.New("not providing password")
	}

	return username, password, nil
}

func (a *Auth) AuthenticateMW() muxer.MidFunc {
	return func(next muxer.HandlerFunc) muxer.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			signed, err := parseBearer(r.Header.Get("authorization"))
			if err != nil {
				return response.WriteError(w, http.StatusUnauthorized, err)
			}

			claims, err := a.parseToken(signed)
			if err != nil {
				return response.WriteError(w, http.StatusUnauthorized, err)
			}

			if err := a.verifyClaims(ctx, claims); err != nil {
				return response.WriteError(w, http.StatusUnauthorized, err)
			}

			ctx = setClaims(ctx, claims)

			return next(ctx, w, r)
		}
	}
}

func parseBearer(bearerToken string) (string, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", errors.New("expected authorization header format: Bearer <token>")
	}

	signed := bearerToken[7:]

	return signed, nil
}

var ErrUnauthorized = errors.New("unauthorized, missing permission")

func (a *Auth) AuthorizeMW(allowedRoles ...string) muxer.MidFunc {
	return func(next muxer.HandlerFunc) muxer.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, err := GetClaims(ctx)
			if err != nil {
				return response.WriteError(w, http.StatusInternalServerError, err)
			}

			if slices.Contains(allowedRoles, claims.Role) {
				return next(ctx, w, r)
			}

			return response.WriteError(w, http.StatusUnauthorized, ErrUnauthorized)
		}
	}
}
