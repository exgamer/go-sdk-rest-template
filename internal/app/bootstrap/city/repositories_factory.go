package city

import (
	cityhttp "github.com/exgamer/go-sdk-rest-template/internal/infrastructure/http/handbook/city"
	citypostgres "github.com/exgamer/go-sdk-rest-template/internal/infrastructure/postgres/handbook/city"
	cityredis "github.com/exgamer/go-sdk-rest-template/internal/infrastructure/redis/handbbok/city"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// redisClient может быть nil, если RedisKernel не зарегистрирован в
// приложении - в этом случае CacheRepository останется nil и Service
// будет работать напрямую с БД (см. errorreporter.CaptureSoft в Service.GetById).
func newRepositoriesFactory(client *gorm.DB, redisClient *redis.Client) *repositoriesFactory {
	factory := &repositoriesFactory{
		PostgresRepository: citypostgres.NewPostgresRepository(client),
		HttpRepository:     cityhttp.NewHttpRepository(),
	}

	if redisClient != nil {
		factory.CacheRepository = cityredis.NewRedisRepository(redisClient)
	}

	return factory
}

type repositoriesFactory struct {
	PostgresRepository *citypostgres.PostgresRepository
	HttpRepository     *cityhttp.HttpRepository
	CacheRepository    *cityredis.RedisRepository
}
