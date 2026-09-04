// Package console собирает консольное (ops) приложение — по аналогии с
// internal/app (основной сервис), только для cmd/console вместо main.go в
// корне. В отличие от internal/app.NewApp не поднимает HTTP/Rabbit kernels:
// консоль не сервис, а разовая команда, выполняется и завершается.
package console

import (
	"github.com/exgamer/go-sdk-rest-template/internal/migrations"
	consolekernel "github.com/exgamer/gosdk-console-core/pkg/app"
	"github.com/exgamer/gosdk-core/pkg/app"
	"github.com/exgamer/gosdk-core/pkg/config"
	postgres "github.com/exgamer/gosdk-postgres-core/pkg/app"
	database "github.com/exgamer/gosdk-postgres-core/pkg/helpers"
	sentryapp "github.com/exgamer/gosdk-sentry-core/pkg/app"
	"gorm.io/gorm"
)

type App struct {
	*app.App
	db *gorm.DB
}

// NewApp открывает соединение с БД и регистрирует все группы ops-команд
// (сейчас только migrate) в ConsoleKernel.
func NewApp() (*App, error) {
	if _, err := config.LoadEnv(); err != nil {
		return nil, err
	}

	db, err := openDefaultConnection()
	if err != nil {
		return nil, err
	}

	migrator := postgres.NewMigrator(db, migrations.All()...)

	consoleKernel := consolekernel.NewConsoleKernel("console").
		AddCommand(migrateCmd(migrator))

	appInstance := &App{
		App: app.NewApp(),
		db:  db,
	}

	if err := appInstance.RegisterAndInitKernels(&sentryapp.SentryKernel{}, consoleKernel); err != nil {
		return nil, err
	}

	return appInstance, nil
}

// Close закрывает соединение с БД. ConsoleKernel не долгоживущий —
// WaitForShutdown() (а вместе с ней и обычный kernel Stop() через stop-hook)
// тут не вызывается, поэтому закрываем явно.
func (a *App) Close() {
	closeConnection(a.db)
}

func openDefaultConnection() (*gorm.DB, error) {
	dbConfig, err := database.InitPostgresDbConfig()
	if err != nil {
		return nil, err
	}

	return database.InitPostgresGormConnection(dbConfig)
}

func closeConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	_ = sqlDB.Close()
}
