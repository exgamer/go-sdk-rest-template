package responses

import "github.com/exgamer/gosdk-http-core/pkg/structures"

type OneItemResponse struct {
	structures.Response[CityItem]
}
