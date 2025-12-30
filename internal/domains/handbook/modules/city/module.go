package city

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/factories"
	"github.com/exgamer/gosdk-core/pkg/app"
	postgresApp "github.com/exgamer/gosdk-postgres-core/pkg/app"
)

// Module модуль городов
type Module struct {
}

func (m *Module) Name() string {
	return "city"
}

func (m *Module) Init(a *app.App) error {
	client, err := postgresApp.GetDefaultPostgresConnection(a)
	if err != nil {
		return err
	}

	repositoryFactory := factories.NewCityRepositoriesFactory(client)
	httpRepositoryFactory := factories.NewCityHttpRepositoriesFactory()
	servicesFactory := factories.NewCityServicesFactory(repositoryFactory, httpRepositoryFactory)
	handlersFactory := factories.NewCityHandlersFactory(servicesFactory)

	err = SetRoutes(a, handlersFactory)

	if err != nil {
		return err
	}

	return nil
}
