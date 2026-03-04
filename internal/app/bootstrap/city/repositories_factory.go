package city

import (
	cityhttp "github.com/exgamer/go-sdk-rest-template/internal/infrastructure/httpclient/handbook/city"
	citypostgres "github.com/exgamer/go-sdk-rest-template/internal/infrastructure/postgres/handbook/city"
	"gorm.io/gorm"
)

func newRepositoriesFactory(client *gorm.DB) *repositoriesFactory {
	return &repositoriesFactory{
		PostgresRepository: citypostgres.NewPostgresRepository(client),
		HttpRepository:     cityhttp.NewHttpRepository(),
	}
}

type repositoriesFactory struct {
	PostgresRepository *citypostgres.PostgresRepository
	HttpRepository     *cityhttp.HttpRepository
}
