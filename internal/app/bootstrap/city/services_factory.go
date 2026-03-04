package city

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"
)

func newServicesFactory(
	repositoryFactory *repositoriesFactory,
) *servicesFactory {
	return &servicesFactory{
		CityService: city.NewService(repositoryFactory.PostgresRepository),
	}
}

type servicesFactory struct {
	CityService *city.Service
}
