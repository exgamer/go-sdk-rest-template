package factories

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/http/city/v1/handlers"
)

func NewCityHandlersFactory(cityServicesFactory *CityServicesFactory) *CityHandlersFactory {
	return &CityHandlersFactory{
		CityHandler: handlers.NewCityHandler(cityServicesFactory.CityService),
	}
}

type CityHandlersFactory struct {
	CityHandler *handlers.CityHandler
}
