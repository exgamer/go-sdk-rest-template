package app

import (
	"github.com/exgamer/go-sdk-rest-template/internal/app/bootstrap/city"
	"github.com/exgamer/gosdk-core/pkg/app"
	http "github.com/exgamer/gosdk-http-core/pkg/app"
	postgres "github.com/exgamer/gosdk-postgres-core/pkg/app"
	rabbitapp "github.com/exgamer/gosdk-rabbit-core/pkg/app"
	sentryapp "github.com/exgamer/gosdk-sentry-core/pkg/app"
)

type App struct {
	*app.App
}

func NewApp() (*App, error) {
	appInstance := &App{
		App: app.NewApp(),
	}

	err := appInstance.RegisterAndInitKernels(
		&sentryapp.SentryKernel{}, // репортинг ошибок в Sentry - errorreporter.Capture/CaptureSoft/CaptureError
		&postgres.PostgresKernel{},
		&http.HttpKernel{},
		rabbitapp.NewRabbitKernel().EnableConsumer().EnablePublisher(), // просто для примера работы с несколькими ядрами которые запускаются
		// &redisapp.RedisKernel{}, // опционально: кеш для city.Service (см. CacheRepository в internal/domains/handbook/city)
	)
	if err != nil {
		return nil, err
	}

	err = appInstance.RegisterAndInitModules(
		&city.Module{},
	)
	if err != nil {
		return nil, err
	}

	return appInstance, nil
}
