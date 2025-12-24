package factories

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/repositories"
)

func NewCityRepositoriesFactory() *CityRepositoriesFactory {
	return &CityRepositoriesFactory{
		CityRepository: repositories.NewCityRepository(),
	}
}

type CityRepositoriesFactory struct {
	CityRepository *repositories.CityRepository
}
