package factories

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/database/dto"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/database/models"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/http/city/v1/requests"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/http/city/v1/responses"
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

func CitySearch(req requests.CityIndexRequest) *dto.CitySearch {
	return &dto.CitySearch{
		Id:   req.Id,
		Name: req.Name,
	}
}

func CityModelFromCreateRequest(req requests.CityCreateRequest) *models.City {
	return &models.City{
		Name: req.Name,
	}
}

func CityModelFromUpdateRequest(req requests.CityUpdateRequest) *models.City {
	return &models.City{
		Name: req.Name,
	}
}

func OneResponse(item *models.City) *responses.CityItem {
	return &responses.CityItem{
		Id:   item.ID,
		Name: item.Name,
	}
}

func PaginatedResponse(paginated *pagination.Paginated[models.City]) *responses.CityPaginated {
	e := &responses.CityPaginated{}
	e.Pagination = *paginated.Pagination
	items := make([]*responses.CityItem, 0)

	for _, item := range paginated.Items {
		items = append(items, OneResponse(item))
	}

	e.Items = items

	return e
}
