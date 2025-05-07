package userapi

import (
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/usercore"
)

type UserResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Role    string    `json:"role"`
	Profile string    `json:"profile"`
}

func toUserResponse(user usercore.User) UserResponse {
	return UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Role:    user.Role,
		Profile: user.Profile,
	}
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func toCreateUserParams(userReq CreateUserRequest) usercore.CreateUserParams {
	return usercore.CreateUserParams{
		Name:     userReq.Name,
		Email:    userReq.Email,
		Password: userReq.Password,
	}
}
