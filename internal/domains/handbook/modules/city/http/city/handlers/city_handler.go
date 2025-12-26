package handlers

import (
	"errors"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/http/city/responses"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/models"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/services"
	"github.com/exgamer/gosdk-http-core/pkg/helpers"
	_ "github.com/exgamer/gosdk-http-core/pkg/structures"
	"github.com/exgamer/gosdk-http-core/validators"
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
//		        @Param			id	path		int	true	"ID города"
//				@Success		200 {object}   responses.OneItemResponse
//				@Failure        404  {object}  structures.NotFoundErrorResponse   "Запись не найдена"
//				@Failure        500  {object}  structures.InternalServerResponse  "Внутренняя ошибка сервера"
//				@Router			/rest-template/v1/city/{id} [get]
func (h *CityHandler) View() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.GetIntQueryParam(c, "id")

		if err != nil {
			return
		}

		item, err := h.cityService.GetCity(c.Request.Context(), uint(id))

		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err, nil)

			return
		}

		if item == nil {
			helpers.ErrorResponse(c, http.StatusNotFound, errors.New("Not Found"), nil)

			return
		}

		helpers.SuccessResponse(c, h.oneResponse(item))
	}
}

func (h *CityHandler) oneResponse(item *models.City) *responses.CityItem {
	return &responses.CityItem{
		Name: item.Name,
	}
}
