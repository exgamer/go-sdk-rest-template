package city

import (
	"github.com/exgamer/gosdk-core/pkg/app"
	"github.com/exgamer/gosdk-http-core/pkg/di"
	"github.com/exgamer/gosdk-http-core/pkg/middleware"
)

func SetRoutes(
	a *app.App,
	handler *Handler,
) error {
	router, err := di.GetRouter(a.Container)
	if err != nil {
		return err
	}

	service := router.Group("/rest-template")
	{
		// Middleware
		service.Use(middleware.RequestInfoMiddleware(a)) // заполнение структур по инфе базовый и http
		service.Use(middleware.LoggerMiddleware())       // форматированные логи
		service.Use(middleware.DebugMiddleware())        // дебаг инфа в ответе от сервиса (только для DEBUG=true)
		service.Use(middleware.SentryMiddleware())       // мидлвейр для отправки ошибок в сентри (если указан SENTRY_DSN= и AppException.TrackInSentry = true)

		v1 := service.Group("/v1")
		{
			v1.Use(middleware.FormattedResponseMiddleware()) // форматированный ответ
			v1.Use(middleware.MetricsMiddleware(a))          // метрики прометея
			v1.GET("/cities", handler.Index())
			v1.GET("/city/:id", handler.View())
			v1.POST("/city", handler.Create())
			v1.PUT("/city/:id", handler.Update())
			v1.DELETE("/city/:id", handler.Delete())
			v1.GET("/city/by-http", handler.ViewByHttp())
		}
	}

	return nil
}
