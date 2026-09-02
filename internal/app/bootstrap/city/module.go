package city

import (
	cityadmin "github.com/exgamer/go-sdk-rest-template/internal/entrypoint/admin/http/handbook/city"
	cityconsumer "github.com/exgamer/go-sdk-rest-template/internal/entrypoint/base/rabbit/handbbok/city"
	"github.com/exgamer/gosdk-core/pkg/app"
	"github.com/exgamer/gosdk-core/pkg/logger"
	"github.com/exgamer/gosdk-postgres-core/pkg/di"
	rDi "github.com/exgamer/gosdk-rabbit-core/pkg/di"
	redisDi "github.com/exgamer/gosdk-redis-core/pkg/di"
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

	// Кеш опционален: если RedisKernel не зарегистрирован в приложении
	// (см. internal/app/app.go), просто работаем без кеша - это не ошибка
	// инициализации модуля.
	redisClient, err := redisDi.GetRedisClient(a.Container)
	if err != nil {
		logger.Info(a.GetContext(), "city: redis cache disabled (RedisKernel not registered)")
		redisClient = nil
	}

	repositoryFactory := newRepositoriesFactory(client, redisClient)
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
