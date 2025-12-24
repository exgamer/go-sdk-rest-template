package factories

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/services"
)

func NewCityServicesFactory(
	repositoryFactory *CityRepositoriesFactory,
) *CityServicesFactory {
	return &CityServicesFactory{
		CityService: services.NewCityService(repositoryFactory.CityRepository),
	}
}

type CityServicesFactory struct {
	CityService *services.CityService
}
