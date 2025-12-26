package factories

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/database/dto"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/http/city/v1/requests"
)

func CitySearchFromIndexRequest(req requests.CityIndexRequest) *dto.CitySearch {
	return &dto.CitySearch{
		Id:   req.Id,
		Name: req.Name,
	}
}
