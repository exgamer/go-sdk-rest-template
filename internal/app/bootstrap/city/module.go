package city

import (
	cityadmin "github.com/exgamer/go-sdk-rest-template/internal/contexts/admin/http/handbook/city"
	cityconsumer "github.com/exgamer/go-sdk-rest-template/internal/contexts/consumer/handbbok/city"
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

	repositoryFactory := newRepositoriesFactory(client)
	servicesFactory := newServicesFactory(repositoryFactory)
	handlersFactory := newHandlersFactory(servicesFactory)

	err = cityadmin.SetRoutes(a, handlersFactory.CityHandler)

	if err != nil {
		return err
	}

	//регистрируем консьюмеры
	consumersFactory := newConsumersFactory()

	consumers := cityconsumer.GetConsumers(consumersFactory.CityConsumer)

	reg, err := rDi.GetRabbitConsumersRegistry(a.Container) // твой helper для DI
	if err != nil {
		return err
	}

	reg.RegisterMultipleHandler(consumers)

	return nil
}
