package city

import (
	"context"
	citydomain "github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"
	"github.com/exgamer/gosdk-http-core/pkg/constants"
	"github.com/exgamer/gosdk-http-core/pkg/gin"
	"github.com/exgamer/gosdk-http-request-builder/pkg/builder"
)

func NewHttpRepository() *HttpRepository {
	return &HttpRepository{}
}

type HttpRepository struct {
}

func (r *HttpRepository) GetCity(ctx context.Context) (*citydomain.City, error) {
	httpInfo := gin.GetHttpInfoFromContext(ctx)
	b := builder.NewGetHttpRequestBuilder[builder.Response[city]](ctx, "http://0.0.0.0:8090/rest-template/v1/city/1").
		SetRequestHeaders(map[string]string{
			constants.RequestIdHeaderName: httpInfo.RequestId,
			"Example-Header":              "example-header-value",
		})

	resp, err := b.GetResult()
	if err != nil {
		return nil, err
	}

	entity := modelToEntity(&resp.Result.Data)

	return entity, nil
}
