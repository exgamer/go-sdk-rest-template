package app

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city"
	"github.com/exgamer/gosdk-core/pkg/app"
	http "github.com/exgamer/gosdk-http-core/pkg/app"
	postgresCore "github.com/exgamer/gosdk-postgres-core/pkg/app"
	app2 "github.com/exgamer/gosdk-rabbit-core/pkg/app"
)

type App struct {
	*app.App
}

func NewApp() (*App, error) {
	appInstance := &App{
		App: app.NewApp(),
	}

	err := appInstance.RegisterAndInitKernels(
		&postgresCore.PostgresKernel{},
		&http.HttpKernel{},
		app2.NewRabbitKernel().EnableConsumer().EnablePublisher(), // просто для примера работы с несколькими ядрами которые запускаются
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

func (a *App) RunHttpKernel() error {
	if err := a.RunKernel(http.HttpKernelName); err != nil {
		return err
	}

	return nil
}

func (a *App) RunRabbitKernel() error {
	if err := a.RunKernel(app2.RabbitKernelName); err != nil {
		return err
	}

	return nil
}
