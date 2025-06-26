package cityapi

import (
	"github.com/google/uuid"
	"github.com/hemozeetah/journi/internal/domain/citycore"
)

type CityResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Caption string    `json:"caption"`
	Images  []string  `json:"imagesURL"`
}

func toCityResponse(city citycore.City) CityResponse {
	images := make([]string, len(city.Images))
	for i, v := range city.Images {
		images[i] = "/static/" + v
	}

	return CityResponse{
		ID:      city.ID,
		Name:    city.Name,
		Caption: city.Caption,
		Images:  images,
	}
}

type CreateCityRequest struct {
	Name    string `form:"name" validate:"required"`
	Caption string `form:"caption" validate:"required"`
}

func toCreateCityParams(cityReq CreateCityRequest, images []string) citycore.CreateCityParams {
	return citycore.CreateCityParams{
		Name:    cityReq.Name,
		Caption: cityReq.Caption,
		Images:  images,
	}
}

type UpdateCityRequest struct {
	Name    *string `json:"name" validate:"omitempty,required"`
	Caption *string `json:"caption" validate:"omitempty,required"`
}

func toUpdateCityParams(cityReq UpdateCityRequest, images []string) citycore.UpdateCityParams {
	var imgs *[]string
	if len(images) != 0 {
		imgs = &images
	}

	return citycore.UpdateCityParams{
		Name:    cityReq.Name,
		Caption: cityReq.Caption,
		Images:  imgs,
	}
}
