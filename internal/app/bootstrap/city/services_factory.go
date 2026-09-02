package city

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"
)

func newServicesFactory(
	repositoryFactory *repositoriesFactory,
) *servicesFactory {
	// ВАЖНО: nil-проверку делаем на конкретном типе *cityredis.RedisRepository,
	// а не после присвоения в интерфейс city.CacheRepository - иначе получим
	// классическую ловушку Go: интерфейс с нетипизированным nil внутри уже
	// не равен nil сам по себе.
	var cacheRepository city.CacheRepository
	if repositoryFactory.CacheRepository != nil {
		cacheRepository = repositoryFactory.CacheRepository
	}

	return &servicesFactory{
		CityService: city.NewService(repositoryFactory.PostgresRepository, cacheRepository),
	}
}

type servicesFactory struct {
	CityService *city.Service
}
