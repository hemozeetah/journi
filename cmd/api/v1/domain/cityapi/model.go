package cityapi

import (
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/citycore"
)

type CityResponse struct {
	ID      uuid.UUID
	Name    string
	Caption string
	Images  []string
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
	Name    string
	Caption string
}

func toCreateCityParams(cityReq CreateCityRequest) citycore.CreateCityParams {
	return citycore.CreateCityParams{
		Name:    cityReq.Name,
		Caption: cityReq.Caption,
	}
}

type UpdateCityRequest struct {
	Name    *string
	Caption *string
}

func toUpdateCityParams(cityReq UpdateCityRequest) citycore.UpdateCityParams {
	return citycore.UpdateCityParams{
		Name:    cityReq.Name,
		Caption: cityReq.Caption,
	}
}
