package postapi

import (
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/cmd/api/v1/jwtauth"
	"github.com/hemozeetah/journi/internal/domain/postcore"
)

type PostResponse struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"userID"`
	PlaceID uuid.UUID `json:"placeID"`
	Caption string    `json:"caption"`
	Images  []string  `json:"images"`
}

func toPostResponse(post postcore.Post) PostResponse {
	return PostResponse{
		ID:      post.ID,
		UserID:  post.UserID,
		PlaceID: post.PlaceID,
		Caption: post.Caption,
		Images:  post.Images,
	}
}

type CreatePostRequest struct {
	PlaceID uuid.UUID `json:"placeID" validate:"required"`
	Caption string    `json:"caption" validate:"required"`
}

func toCreatePostParams(postReq CreatePostRequest, claims jwtauth.Claims) postcore.CreatePostParams {
	return postcore.CreatePostParams{
		UserID:  claims.ID,
		PlaceID: postReq.PlaceID,
		Caption: postReq.Caption,
	}
}

type UpdatePostRequest struct {
	PlaceID *uuid.UUID `json:"placeID" validate:"omitempty,required"`
	Caption *string    `json:"caption" validate:"omitempty,required"`
}

func toUpdatePostParams(postReq UpdatePostRequest) postcore.UpdatePostParams {
	return postcore.UpdatePostParams{
		PlaceID: postReq.PlaceID,
		Caption: postReq.Caption,
	}
}
