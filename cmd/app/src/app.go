package src

import (
	"database/sql"
	"exgamer.kz/example-service/cmd/app/src/factories"
	"exgamer.kz/example-service/cmd/app/src/routes"
	structures2 "github.com/exgamer/go-rest-sdk/pkg/config/structures"
	helpers "github.com/exgamer/go-rest-sdk/pkg/helpers/db/mysql"
	"github.com/exgamer/go-rest-sdk/pkg/helpers/db/redis"
	ginhelper "github.com/exgamer/go-rest-sdk/pkg/helpers/gin"
	"github.com/gin-gonic/gin"
	redis2 "github.com/redis/go-redis/v9"
	"log"
)

func NewApp(
	appConfig *structures2.AppConfig,
	dbConfig *structures2.DbConfig,
	redisConfig *structures2.RedisConfig,
) *App {
	return &App{
		appConfig:   appConfig,
		dbConfig:    dbConfig,
		redisConfig: redisConfig,
	}
}

type App struct {
	appConfig             *structures2.AppConfig
	dbConfig              *structures2.DbConfig
	redisConfig           *structures2.RedisConfig
	db                    *sql.DB
	mysqlRepoFactory      *factories.MysqlRepositoryFactory
	entityManagersFactory *factories.EntityManagersFactory
	servicesFactory       *factories.ServicesFactory
	controllersFactory    *factories.ControllersFactory
	router                *gin.Engine
	redisClient           *redis2.Client
}

func (app *App) Run() {
	app.OpenDbConnection()
	defer app.CloseDbConnection()

	app.OpenRedisConnection()
	defer app.CloseRedisConnection()

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

func (app *App) OpenRedisConnection() {
	// Open up database connection.
	client := redis.OpenRedisConnection(app.redisConfig)

	app.redisClient = client
}

func (app *App) CloseRedisConnection() {
	redis.CloseRedisConnection(app.redisClient)
}
