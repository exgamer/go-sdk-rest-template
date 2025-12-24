package city

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/factories"
	"github.com/exgamer/gosdk-core/pkg/app"
)

// Module модуль городов
type Module struct {
}

func (m *Module) Name() string {
	return "city"
}

func (m *Module) Init(a *app.App) error {
	repositoryFactory := factories.NewCityRepositoriesFactory()
	servicesFactory := factories.NewCityServicesFactory(repositoryFactory)
	handlersFactory := factories.NewCityHandlersFactory(servicesFactory)

	err := SetRoutes(a, handlersFactory)

	if err != nil {
		return err
	}

	return nil
}
