package handlers

import (
	"errors"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/http/city/v1/factories"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/http/city/v1/requests"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/services"
	"github.com/exgamer/gosdk-http-core/pkg/helpers"
	"github.com/exgamer/gosdk-http-core/pkg/response"
	_ "github.com/exgamer/gosdk-http-core/pkg/structures"
	"github.com/exgamer/gosdk-http-core/pkg/validators"
	"github.com/gin-gonic/gin"
)

func NewCityHandler(cityService *services.CityService, cityHttpService *services.CityHttpService) *CityHandler {
	return &CityHandler{
		cityService:     cityService,
		cityHttpService: cityHttpService,
	}
}

type CityHandler struct {
	cityService     *services.CityService
	cityHttpService *services.CityHttpService
}

//		@Summary		Список городов
//		@Description	Список городов
//		@Tags			city
//		@Accept			json
//		@Produce		json
//	    @Param			message	query		requests.CityIndexRequest	true	"Search Cities"
//		@Success		200 {object} responses.CityPaginatedResponse
//	    @Failure        500  {object} structures.InternalServerResponse  "Внутренняя ошибка сервера"
//		@Router			/rest-template/v1/cities [get]
func (h *CityHandler) Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		pagerRequest, err := helpers.GetPagerRequest(c)

		if err != nil {
			response.UnprocessableEntity(c, err, nil)

			return
		}

		request := requests.CityIndexRequest{}
		ve := validators.ValidateRequestQuery(c, &request)

		if ve == false {
			return
		}

		searchDto := factories.CitySearch(request)
		searchDto.Page = pagerRequest.Page
		searchDto.PerPage = pagerRequest.PerPage
		paginated, err := h.cityService.Paginated(c.Request.Context(), searchDto)

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		response.Success(c, factories.PaginatedResponse(paginated))
	}
}

//			@Summary		Возвращает город по id
//			@Description	Возвращает город по id
//			@Tags			city
//			@Accept			json
//			@Produce		json
//	        @Param			id	path		int	true	"ID города"
//			@Success		200 {object}   responses.CityItemResponse
//		    @Failure		400	{object}	structures.BadRequestErrorResponse				"Плохой запрос"
//			@Failure        404  {object}  structures.NotFoundErrorResponse   "Запись не найдена"
//			@Failure        500  {object}  structures.InternalServerResponse  "Внутренняя ошибка сервера"
//			@Router			/rest-template/v1/city/{id} [get]
func (h *CityHandler) View() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.GetIntQueryParam(c, "id")

		if err != nil {
			response.BadRequest(c, err, nil)

			return
		}

		item, err := h.cityService.GetById(c.Request.Context(), uint(id))

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		if item == nil {
			response.NotFound(c, errors.New("Not Found"), nil)

			return
		}

		response.Success(c, factories.OneResponse(item))
	}
}

//			@Summary		Создать город
//			@Description	Создать город
//			@Tags			city
//			@Accept			json
//			@Produce		json
//			@Param			message	body requests.CityCreateRequest	true	"Create city"
//			@Success		201 {object} responses.CityItemResponse
//		    @Failure		422	{object}	structures.ValidationErrorResponse				"Ошибка валидации параметров"
//		    @Failure		500	{object}	structures.InternalServerResponse				"Внутренняя ошибка сервера"
//	     @Router			/rest-template/v1/city/ [post]
func (h *CityHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := requests.CityCreateRequest{}
		ve := validators.ValidateRequestBody(c, &request)

		if ve == false {
			return
		}

		m := factories.CityModelFromCreateRequest(request)
		model, err := h.cityService.Create(c.Request.Context(), m)

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		response.SuccessCreated(c, factories.OneResponse(model))

		return
	}
}

//			@Summary		Создать город
//			@Description	Создать город
//			@Tags			city
//			@Accept			json
//			@Produce		json
//		    @Param			id	path		int	true	"ID города"
//			@Param			message	body requests.CityCreateRequest	true	"Create city"
//			@Success		201 {object} responses.CityItemResponse
//		    @Failure		400	{object}	structures.BadRequestErrorResponse				"Плохой запрос"
//		    @Failure		422	{object}	structures.ValidationErrorResponse				"Ошибка валидации параметров"
//		    @Failure		500	{object}	structures.InternalServerResponse				"Внутренняя ошибка сервера"
//	     @Router			/rest-template/v1/city/{id} [put]
func (h *CityHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.GetIntQueryParam(c, "id")

		if err != nil {
			response.BadRequest(c, err, nil)

			return
		}

		request := requests.CityCreateRequest{}
		ve := validators.ValidateRequestBody(c, &request)

		if ve == false {
			return
		}

		m := factories.CityModelFromCreateRequest(request)
		m.ID = uint(id)
		model, err := h.cityService.Update(c.Request.Context(), m)

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		response.Success(c, factories.OneResponse(model))

		return
	}
}

//			@Summary		удалить город
//			@Description	удалить город
//			@Tags			city
//			@Accept			json
//			@Produce		json
//		    @Success		204
//			Failure		    400	{object}	structures.BadRequestErrorResponse				"Плохой запрос"
//	     @Failure		500	{object}	structures.InternalServerResponse				"Внутренняя ошибка сервера"
//	     @Router			/rest-template/v1/city/{id} [delete]
func (h *CityHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.GetIntQueryParam(c, "id")

		if err != nil {
			response.BadRequest(c, err, nil)

			return
		}

		err = h.cityService.Delete(c.Request.Context(), uint(id))

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		response.SuccessDeleted(c, nil)

		return
	}
}

func (h *CityHandler) ViewByHttp() gin.HandlerFunc {
	return func(c *gin.Context) {
		item, err := h.cityHttpService.GetCity(c.Request.Context())

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		response.Success(c, factories.OneResponseFromDto(item))
	}
}
