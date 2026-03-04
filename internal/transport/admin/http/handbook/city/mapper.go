package city

import (
	city2 "github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

func searchFromIndexRequest(req indexRequest) *city2.Search {
	return &city2.Search{
		Id:   req.Id,
		Name: req.Name,
	}
}

func entityFromCreateRequest(req createRequest) *city2.City {
	return &city2.City{
		Name: req.Name,
	}
}

func entityFromUpdateRequest(req updateRequest) *city2.City {
	return &city2.City{
		Name: req.Name,
	}
}

//func itemFromEntity(item *dto2.City) *Item {
//	return &Item{
//		Id:   item.ID,
//		Name: item.Name,
//	}
//}

func itemFromEntity(city *city2.City) *item {
	return &item{
		Id:   city.ID,
		Name: city.Name,
	}
}

func paginatedResponseFromPagination(paginated *pagination.Paginated[city2.City]) *paginatedItem {
	e := &paginatedItem{}
	e.Pagination = *paginated.Pagination
	items := make([]*item, 0)

	for _, item := range paginated.Items {
		items = append(items, itemFromEntity(item))
	}

	e.Items = items

	return e
}
