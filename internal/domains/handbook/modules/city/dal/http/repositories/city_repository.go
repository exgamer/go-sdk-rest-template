package repositories

import (
	"context"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/http/dto"
	"github.com/exgamer/gosdk-http-core/pkg/constants"
	"github.com/exgamer/gosdk-http-core/pkg/gin"
	"github.com/exgamer/gosdk-http-request-builder/pkg/builder"
)

func NewCityRepository() *CityRepository {
	return &CityRepository{}
}

type CityRepository struct {
}

func (r *CityRepository) GetCity(ctx context.Context) (*dto.City, error) {
	httpInfo := gin.GetHttpInfoFromContext(ctx)
	b := builder.NewGetHttpRequestBuilder[dto.City](ctx, "http://0.0.0.0:8090/rest-template/v1/city/1").
		SetRequestHeaders(map[string]string{
			constants.RequestIdHeaderName: httpInfo.RequestId,
			"Example-Header":              "example-header-value",
		})

	resp, err := b.GetResult()
	if err != nil {
		return nil, err
	}

	return &resp.Result.Data, nil
}
