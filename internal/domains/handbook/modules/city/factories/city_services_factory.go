package factories

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/services"
)

func NewCityServicesFactory(
	repositoryFactory *CityRepositoriesFactory,
	cityHttpRepositoriesFactory *CityHttpRepositoriesFactory,
) *CityServicesFactory {
	return &CityServicesFactory{
		CityService:     services.NewCityService(repositoryFactory.CityRepository),
		CityHttpService: services.NewCityHttpService(cityHttpRepositoriesFactory.CityRepository),
	}
}

type CityServicesFactory struct {
	CityService     *services.CityService
	CityHttpService *services.CityHttpService
}
