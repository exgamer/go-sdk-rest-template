package factories

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/http/repositories"
)

func NewCityHttpRepositoriesFactory() *CityHttpRepositoriesFactory {
	return &CityHttpRepositoriesFactory{
		CityRepository: repositories.NewCityRepository(),
	}
}

type CityHttpRepositoriesFactory struct {
	CityRepository *repositories.CityRepository
}
