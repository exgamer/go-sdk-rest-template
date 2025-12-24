package handlers

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/http/city/responses"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/models"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/services"
	"github.com/exgamer/gosdk-http-core/pkg/helpers"
	_ "github.com/exgamer/gosdk-http-core/pkg/structures"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewCityHandler(cityService *services.CityService) *CityHandler {
	return &CityHandler{
		cityService: cityService,
	}
}

type CityHandler struct {
	cityService *services.CityService
}

//	         Возвращает город Алматы
//
//				@Summary		Возвращает город Алматы
//				@Description	Возвращает город Алматы
//				@Tags			city
//				@Accept			json
//				@Produce		json
//				@Success		200 {object} responses.CityItemResponse
//				@Failure        500  {object}  structures.InternalServerResponse  "Внутренняя ошибка сервера"
//				@Router			/rest-template/v1/city [get]
func (h *CityHandler) GetCity() gin.HandlerFunc {
	return func(c *gin.Context) {
		item, err := h.cityService.GetCity()

		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err, nil)

			return
		}

		helpers.SuccessResponse(c, h.oneResponse(item))
	}
}

func (h *CityHandler) oneResponse(item *models.City) *responses.CityItemResponse {
	return &responses.CityItemResponse{
		Name: item.Name,
	}
}
