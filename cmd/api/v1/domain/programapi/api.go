package programapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/cmd/api/v1/request"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/programcore"
	"github.com/hemozeetah/journi/pkg/logger"
)

type api struct {
	log  *logger.Logger
	core *programcore.Core
}

func (a *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var programReq CreateProgramRequest
	if err := request.ParseBody(r, &programReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(programReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	claims, err := jwtauth.GetClaims(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	program, err := a.core.Create(ctx, toCreateProgramParams(programReq, claims))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	programResp := toProgramResponse(program)
	return response.Write(w, http.StatusCreated, programResp)
}

func (a *api) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	program, err := getProgram(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	programResp := toProgramResponse(program)
	return response.Write(w, http.StatusOK, programResp)
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
	programs, err := a.core.Query(ctx, query)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	programsResp := make([]ProgramResponse, len(programs))
	for i, program := range programs {
		programsResp[i] = toProgramResponse(program)
	}
	return response.Write(w, http.StatusOK, programsResp)
}

func (a *api) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var programReq UpdateProgramRequest
	if err := request.ParseBody(r, &programReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(programReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	program, err := getProgram(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	fmt.Println(programReq)
	program, err = a.core.Update(ctx, program, toUpdateProgramParams(programReq))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	programResp := toProgramResponse(program)
	return response.Write(w, http.StatusOK, programResp)
}

func (a *api) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	program, err := getProgram(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	if err := a.core.Delete(ctx, program); err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	return response.Write(w, http.StatusNoContent, nil)
}
