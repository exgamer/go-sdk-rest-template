package src

import (
	"database/sql"
	"exgamer.kz/example-service/cmd/app/src/factories"
	"exgamer.kz/example-service/cmd/app/src/routes"
	structures2 "github.com/exgamer/go-rest-sdk/pkg/config/structures"
	helpers "github.com/exgamer/go-rest-sdk/pkg/helpers/db/mysql"
	ginhelper "github.com/exgamer/go-rest-sdk/pkg/helpers/gin"
	"github.com/gin-gonic/gin"
	"log"
)

func NewApp(
	appConfig *structures2.AppConfig,
	dbConfig *structures2.DbConfig,
) *App {
	return &App{
		appConfig: appConfig,
		dbConfig:  dbConfig,
	}
}

type App struct {
	appConfig             *structures2.AppConfig
	dbConfig              *structures2.DbConfig
	db                    *sql.DB
	mysqlRepoFactory      *factories.MysqlRepositoryFactory
	entityManagersFactory *factories.EntityManagersFactory
	servicesFactory       *factories.ServicesFactory
	controllersFactory    *factories.ControllersFactory
	router                *gin.Engine
}

func (app *App) Run() {
	app.OpenDbConnection()
	defer app.CloseDbConnection()

	app.router = ginhelper.InitRouter(app.appConfig)

	app.mysqlRepoFactory = factories.NewMysqlRepositoryFactory(app.db)
	app.entityManagersFactory = factories.NewEntityManagersFactory(app.mysqlRepoFactory)
	app.servicesFactory = factories.NewServiceFactory(app.entityManagersFactory)
	app.controllersFactory = factories.NewControllersFactory(app.servicesFactory)

	routes.SetRoutes(app.db, app.router, app.appConfig, app.controllersFactory)
	app.router.Run(app.appConfig.ServerAddress)
}

func (app *App) OpenDbConnection() error {
	// Open up database connection.
	db, err := helpers.OpenMysqlConnection(app.dbConfig)

	if err != nil {
		log.Fatal(err)
	}

	app.db = db

	return nil
}

func (app *App) CloseDbConnection() {
	helpers.CloseMysqlConnection(app.db)
}
