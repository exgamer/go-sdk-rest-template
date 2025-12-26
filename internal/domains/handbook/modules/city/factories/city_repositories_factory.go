package factories

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/repositories"
	"gorm.io/gorm"
)

func NewCityRepositoriesFactory(client *gorm.DB) *CityRepositoriesFactory {
	return &CityRepositoriesFactory{
		CityRepository: repositories.NewCityRepository(client),
	}
}

type CityRepositoriesFactory struct {
	CityRepository *repositories.CityRepository
}
