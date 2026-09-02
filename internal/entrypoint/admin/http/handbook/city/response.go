package city

import (
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
	"github.com/exgamer/gosdk-http-core/pkg/structures"
)

type item struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type itemResponse struct {
	structures.Response[item]
}

type paginatedItem struct {
	Items      []*item               `json:"items"`
	Pagination pagination.Pagination `json:"pagination"`
}

type paginatedResponse struct {
	structures.Response[paginatedItem]
}

// pingResponse пример ответа без обёртки {success, data} — отдаётся через response.RawSuccess.
type pingResponse struct {
	Status string `json:"status"`
}
