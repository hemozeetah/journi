package cityapi

import (
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/citycore"
)

type CityResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Caption string    `json:"caption"`
	Images  []string  `json:"images"`
}

func toCityResponse(city citycore.City) CityResponse {
	return CityResponse{
		ID:      city.ID,
		Name:    city.Name,
		Caption: city.Caption,
		Images:  city.Images,
	}
}

type CreateCityRequest struct {
	Name    string `json:"name" validate:"required"`
	Caption string `json:"caption" validate:"required"`
}

func toCreateCityParams(cityReq CreateCityRequest) citycore.CreateCityParams {
	return citycore.CreateCityParams{
		Name:    cityReq.Name,
		Caption: cityReq.Caption,
	}
}

type UpdateCityRequest struct {
	Name    *string `json:"name" validate:"omitempty,required"`
	Caption *string `json:"caption" validate:"omitempty,required"`
}

func toUpdateCityParams(cityReq UpdateCityRequest) citycore.UpdateCityParams {
	return citycore.UpdateCityParams{
		Name:    cityReq.Name,
		Caption: cityReq.Caption,
	}
}
