package city

import "github.com/exgamer/go-sdk-rest-template/internal/contexts/admin/http/handbook/city"

func newHandlersFactory(cityServicesFactory *servicesFactory) *handlersFactory {
	return &handlersFactory{
		CityHandler: city.NewHandler(cityServicesFactory.CityService),
	}
}

type handlersFactory struct {
	CityHandler *city.Handler
}
