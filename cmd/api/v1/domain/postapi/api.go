package postapi

import (
	"context"
	"net/http"

	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/cmd/api/v1/request"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/postcore"
	"github.com/hemozeetah/journi/pkg/logger"
)

type api struct {
	log  *logger.Logger
	core *postcore.Core
}

func (a *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	var postReq CreatePostRequest
	if err := request.ParseForm(r, &postReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(postReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	images, err := request.ParseFile(r, "images")
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	claims, err := jwtauth.GetClaims(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	post, err := a.core.Create(ctx, toCreatePostParams(postReq, claims, images))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	postResp := toPostResponse(post)
	return response.Write(w, http.StatusCreated, postResp)
}

func (a *api) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	post, err := getPost(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	postResp := toPostResponse(post)
	return response.Write(w, http.StatusOK, postResp)
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

	posts, err := a.core.Query(ctx, query)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	postsResp := make([]PostResponse, len(posts))
	for i, post := range posts {
		postsResp[i] = toPostResponse(post)
	}
	return response.Write(w, http.StatusOK, postsResp)
}

func (a *api) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	var postReq UpdatePostRequest
	if err := request.ParseForm(r, &postReq); err != nil {
		return response.WriteError(w, http.StatusBadRequest, err)
	}

	if err := request.Validate(postReq); err != nil {
		return response.WriteError(w, http.StatusUnprocessableEntity, err)
	}

	images, err := request.ParseFile(r, "images")
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	post, err := getPost(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	post, err = a.core.Update(ctx, post, toUpdatePostParams(postReq, images))
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	postResp := toPostResponse(post)
	return response.Write(w, http.StatusOK, postResp)
}

func (a *api) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	post, err := getPost(ctx)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	if err := a.core.Delete(ctx, post); err != nil {
		return response.WriteError(w, http.StatusInternalServerError, err)
	}

	return response.Write(w, http.StatusNoContent, nil)
}
