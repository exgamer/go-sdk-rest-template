package city

import (
	"errors"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"
	"github.com/exgamer/gosdk-http-core/pkg/helpers"
	"github.com/exgamer/gosdk-http-core/pkg/response"
	_ "github.com/exgamer/gosdk-http-core/pkg/structures"
	"github.com/exgamer/gosdk-http-core/pkg/validators"
	"github.com/gin-gonic/gin"
)

func NewHandler(cityService *city.Service) *Handler {
	return &Handler{
		cityService: cityService,
	}
}

type Handler struct {
	cityService *city.Service
}

//		@Summary		Список городов
//		@Description	Список городов
//		@Tags			city
//		@Accept			json
//		@Produce		json
//	    @Param			message	query		indexRequest	true	"Search Cities"
//		@Success		200 {object} paginatedResponse
//	    @Failure        500  {object} structures.InternalServerResponse  "Внутренняя ошибка сервера"
//		@Router			/rest-template/v1/cities [get]
func (h *Handler) Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		pagerRequest, err := helpers.GetPagerRequest(c)

		if err != nil {
			response.UnprocessableEntity(c, err, nil)

			return
		}

		request := indexRequest{}
		ve := validators.ValidateRequestQuery(c, &request)

		if ve == false {
			return
		}

		searchDto := searchFromIndexRequest(request)
		searchDto.Page = pagerRequest.Page
		searchDto.PerPage = pagerRequest.PerPage
		paginated, err := h.cityService.Paginated(c.Request.Context(), searchDto)

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		response.Success(c, paginatedResponseFromPagination(paginated))
	}
}

//			@Summary		Возвращает город по id
//			@Description	Возвращает город по id
//			@Tags			city
//			@Accept			json
//			@Produce		json
//	        @Param			id	path		int	true	"ID города"
//			@Success		200 {object}   itemResponse
//		    @Failure		400	{object}	structures.BadRequestErrorResponse				"Плохой запрос"
//			@Failure        404  {object}  structures.NotFoundErrorResponse   "Запись не найдена"
//			@Failure        500  {object}  structures.InternalServerResponse  "Внутренняя ошибка сервера"
//			@Router			/rest-template/v1/city/{id} [get]
func (h *Handler) View() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.GetIntQueryParam(c, "id")

		if err != nil {
			response.BadRequest(c, err, nil)

			return
		}

		entity, err := h.cityService.GetById(c.Request.Context(), uint(id))

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		if entity == nil {
			response.NotFound(c, errors.New("Not Found"), nil)

			return
		}

		response.Success(c, itemFromEntity(entity))
	}
}

//			@Summary		Создать город
//			@Description	Создать город
//			@Tags			city
//			@Accept			json
//			@Produce		json
//			@Param			message	body createRequest	true	"Create city"
//			@Success		201 {object} itemResponse
//		    @Failure		422	{object}	structures.ValidationErrorResponse				"Ошибка валидации параметров"
//		    @Failure		500	{object}	structures.InternalServerResponse				"Внутренняя ошибка сервера"
//	     @Router			/rest-template/v1/city/ [post]
func (h *Handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := createRequest{}
		ve := validators.ValidateRequestBody(c, &request)

		if ve == false {
			return
		}

		m := entityFromCreateRequest(request)
		model, err := h.cityService.Create(c.Request.Context(), m)

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		response.SuccessCreated(c, itemFromEntity(model))

		return
	}
}

//			@Summary		Редактировать город
//			@Description	Редактировать город
//			@Tags			city
//			@Accept			json
//			@Produce		json
//		    @Param			id	path		int	true	"ID города"
//			@Param			message	body updateRequest	true	"Create city"
//			@Success		201 {object} itemResponse
//		    @Failure		400	{object}	structures.BadRequestErrorResponse				"Плохой запрос"
//		    @Failure		422	{object}	structures.ValidationErrorResponse				"Ошибка валидации параметров"
//		    @Failure		500	{object}	structures.InternalServerResponse				"Внутренняя ошибка сервера"
//	     @Router			/rest-template/v1/city/{id} [put]
func (h *Handler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := validators.GetIntQueryParam(c, "id")

		if err != nil {
			response.BadRequest(c, err, nil)

			return
		}

		request := updateRequest{}
		ve := validators.ValidateRequestBody(c, &request)

		if ve == false {
			return
		}

		m := entityFromUpdateRequest(request)
		m.ID = uint(id)
		entity, err := h.cityService.Update(c.Request.Context(), m)

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		response.Success(c, itemFromEntity(entity))

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
func (h *Handler) Delete() gin.HandlerFunc {
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

func (h *Handler) ViewByHttp() gin.HandlerFunc {
	return func(c *gin.Context) {
		item, err := h.cityService.GetCity(c.Request.Context())

		if err != nil {
			response.InternalServerError(c, err, nil)

			return
		}

		response.Success(c, itemFromEntity(item))
	}
}
