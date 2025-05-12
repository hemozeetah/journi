package programapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/programcore"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/muxer"
)

const programKey = "program"

func setProgram(ctx context.Context, program programcore.Program) context.Context {
	return context.WithValue(ctx, programKey, program)
}

func getProgram(ctx context.Context) (programcore.Program, error) {
	program, ok := ctx.Value(programKey).(programcore.Program)
	if !ok {
		return programcore.Program{}, errors.New("program not found")
	}

	return program, nil
}

func (a *api) parseProgramMW(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		programID, err := uuid.Parse(r.PathValue("program_id"))
		if err != nil {
			return response.WriteError(w, http.StatusBadRequest, err)
		}

		program, err := a.core.QueryByID(ctx, programID)
		if err != nil {
			if errors.Is(err, programcore.ErrNotFound) {
				return response.WriteError(w, http.StatusNotFound, programcore.ErrNotFound)
			}
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		ctx = setProgram(ctx, program)
		return handler(ctx, w, r)
	}
}

func (a *api) adminOrOwner(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		claims, err := jwtauth.GetClaims(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		program, err := getProgram(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		if claims.Role == usercore.RoleAdmin || claims.ID == program.CompanyID {
			return handler(ctx, w, r)
		}

		return response.WriteError(w, http.StatusUnauthorized, errors.New("not admin or owner"))
	}
}
