package city

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/factories"
	"github.com/exgamer/gosdk-core/pkg/app"
	"github.com/exgamer/gosdk-postgres-core/pkg/di"
	rDi "github.com/exgamer/gosdk-rabbit-core/pkg/di"
)

// Module модуль городов
type Module struct {
}

func (m *Module) Name() string {
	return "city"
}

func (m *Module) Init(a *app.App) error {
	client, err := di.GetDefaultPostgresConnection(a.Container)
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

	//регистрируем консьюмеры
	consumersFactory := NewCityConsumersFactory()

	consumers := GetConsumers(consumersFactory)

	reg, err := rDi.GetRabbitConsumersRegistry(a.Container) // твой helper для DI
	if err != nil {
		return err
	}

	reg.RegisterMultipleHandler(consumers)

	return nil
}
