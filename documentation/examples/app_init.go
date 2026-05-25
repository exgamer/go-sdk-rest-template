//go:build ignore

// Пример инициализации приложения.
// Реальный файл: internal/app/app.go

package app

import (
	"github.com/exgamer/gosdk-core/pkg/app"
	http     "github.com/exgamer/gosdk-http-core/pkg/app"
	postgres "github.com/exgamer/gosdk-postgres-core/pkg/app"
	rabbit   "github.com/exgamer/gosdk-rabbit-core/pkg/app"

	city    "github.com/exgamer/go-sdk-rest-template/internal/app/bootstrap/city"
	product "github.com/exgamer/go-sdk-rest-template/internal/app/bootstrap/product"
)

// ===== internal/app/app.go =====

type App struct {
	*app.App
}

func NewApp() (*App, error) {
	appInstance := &App{App: app.NewApp()}

	err := appInstance.RegisterAndInitKernels(
		&postgres.PostgresKernel{},
		&http.HttpKernel{},
		rabbit.NewRabbitKernel().EnableConsumer().EnablePublisher(),
	)
	if err != nil {
		return nil, err
	}

	err = appInstance.RegisterAndInitModules(
		&city.Module{},
		&product.Module{},
	)
	if err != nil {
		return nil, err
	}

	return appInstance, nil
}

// ===== Доступные варианты Kernels =====

// PostgresKernel — GORM + Postgres
// &postgres.PostgresKernel{}

// HttpKernel — Gin сервер
// &http.HttpKernel{}

// RabbitKernel — RabbitMQ (варианты):
// rabbit.NewRabbitKernel()                               // только соединение
// rabbit.NewRabbitKernel().EnableConsumer()              // + consumer
// rabbit.NewRabbitKernel().EnablePublisher()             // + publisher
// rabbit.NewRabbitKernel().EnableConsumer().EnablePublisher() // оба

// ===== Шаблон Module =====
// Реальный файл: internal/app/bootstrap/{module}/module.go

type Module struct{}

func (m *Module) Name() string { return "module-name" }

func (m *Module) Init(a *app.App) error {
	// 1. получить зависимости из DI
	// 2. создать factory цепочку: repo → service → handler
	// 3. зарегистрировать маршруты и consumers
	return nil
}
