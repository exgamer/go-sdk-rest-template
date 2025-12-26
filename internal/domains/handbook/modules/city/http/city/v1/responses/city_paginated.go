package responses

import (
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

type CityPaginated struct {
	Items      []*CityItem           `json:"items"`
	Pagination pagination.Pagination `json:"pagination"`
}
