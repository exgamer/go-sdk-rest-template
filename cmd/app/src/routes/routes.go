package routes

import (
	"database/sql"
	"exgamer.kz/example-service/cmd/app/src/factories"
	structures2 "github.com/exgamer/go-rest-sdk/pkg/config/structures"
	httpResponse "github.com/exgamer/go-rest-sdk/pkg/http"
	"github.com/exgamer/go-rest-sdk/pkg/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

func SetRoutes(
	db *sql.DB,
	router *gin.Engine,
	appConfig *structures2.AppConfig,
	controllersFactory *factories.ControllersFactory,
) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// healthcheck
	healthcheck := router.Group("/healthcheck")
	{
		healthcheck.Use(middleware.ResponseHandler)
		healthcheck.GET("/readiness", func(c *gin.Context) {
			db.Ping()
			httpResponse.Response(c, http.StatusOK, map[string]any{
				"service_name":   appConfig.Name,
				"pod_name":       appConfig.HostName,
				"date":           time.Now().Format("2006-01-02 15:04:05"),
				"database":       "ok",
				"message_broker": "no broker",
			})
		})
		healthcheck.GET("/liveness", func(c *gin.Context) {
			httpResponse.Response(c, http.StatusOK, map[string]any{
				"service_name": appConfig.Name,
				"pod_name":     appConfig.HostName,
				"date":         time.Now().Format("2006-01-02 15:04:05"),
			})
		})
	}
	// Routers
	v1 := router.Group("/example-go/v1")
	{
		v1.Use(middleware.ResponseHandler)
		v1.Use(middleware.PinbaHandler(appConfig))
		v1.GET("/test", controllersFactory.PostController.All())
	}
}
