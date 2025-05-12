package postapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/cmd/api/v1/response"
	"github.com/hemozeetah/journi/internal/domain/postcore"
	"github.com/hemozeetah/journi/internal/domain/usercore"
	"github.com/hemozeetah/journi/pkg/muxer"
)

const postKey = "post"

func setPost(ctx context.Context, post postcore.Post) context.Context {
	return context.WithValue(ctx, postKey, post)
}

func getPost(ctx context.Context) (postcore.Post, error) {
	post, ok := ctx.Value(postKey).(postcore.Post)
	if !ok {
		return postcore.Post{}, errors.New("post not found")
	}

	return post, nil
}

func (a *api) parsePostMW(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		postID, err := uuid.Parse(r.PathValue("post_id"))
		if err != nil {
			return response.WriteError(w, http.StatusBadRequest, err)
		}

		post, err := a.core.QueryByID(ctx, postID)
		if err != nil {
			if errors.Is(err, postcore.ErrNotFound) {
				return response.WriteError(w, http.StatusNotFound, postcore.ErrNotFound)
			}
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		ctx = setPost(ctx, post)
		return handler(ctx, w, r)
	}
}

func (a *api) adminOrOwner(handler muxer.HandlerFunc) muxer.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		claims, err := jwtauth.GetClaims(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		post, err := getPost(ctx)
		if err != nil {
			return response.WriteError(w, http.StatusInternalServerError, err)
		}

		if claims.Role == usercore.RoleAdmin || claims.ID == post.UserID {
			return handler(ctx, w, r)
		}

		return response.WriteError(w, http.StatusUnauthorized, errors.New("not admin or owner"))
	}
}
