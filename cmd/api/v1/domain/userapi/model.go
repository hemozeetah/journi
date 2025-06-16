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
	Profile string    `json:"profileImageURL"`
}

func toUserResponse(user usercore.User) UserResponse {
	profile := "/static/" + user.Profile

	return UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Role:    user.Role,
		Profile: profile,
	}
}

type CreateUserRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`
}

func toCreateUserParams(userReq CreateUserRequest) usercore.CreateUserParams {
	return usercore.CreateUserParams{
		Name:     userReq.Name,
		Email:    userReq.Email,
		Password: userReq.Password,
	}
}

type UpdateUserRequest struct {
	Name            *string `json:"name" validate:"omitempty,required"`
	Email           *string `json:"email" validate:"omitempty,required,email"`
	Password        *string `json:"password" validate:"omitempty,required"`
	PasswordConfirm *string `json:"passwordConfirm" validate:"required_with=Password,eqfield=Password"`
	Role            *string `json:"role" validate:"omitempty,oneof=user company admin"`
}

func toUpdateUserParams(userReq UpdateUserRequest, images []string) usercore.UpdateUserParams {
	var image *string
	if len(images) != 0 {
		image = &images[0]
	}

	return usercore.UpdateUserParams{
		Name:     userReq.Name,
		Email:    userReq.Email,
		Password: userReq.Password,
		Role:     userReq.Role,
		Profile:  image,
	}
}
