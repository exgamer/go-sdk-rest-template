package postHttp

import (
	"exgamer.kz/example-service/cmd/app/src/modules/post/v1/postServices"
	httpResponse "github.com/exgamer/go-rest-sdk/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewPostController(service *postServices.PostService) *PostController {
	return &PostController{
		postService: service,
	}
}

type PostController struct {
	postService *postServices.PostService
}

// Posts godoc
//	@Summary		Posts
//	@Description	Posts
//	@Tags			gifts
//	@Accept			json
//	@Produce		json
//	@Success		200		{object} postResponses.PostResponse
//	@Router			/example-go/v1/posts [get]
func (controller *PostController) All() gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, appException := controller.postService.All()

		if appException != nil {
			httpResponse.Error(c, appException)

			return
		}

		httpResponse.Response(c, http.StatusOK, posts)
	}
}
