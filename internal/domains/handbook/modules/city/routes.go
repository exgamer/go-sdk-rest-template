package city

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/factories"
	"github.com/exgamer/gosdk-core/pkg/app"
	http "github.com/exgamer/gosdk-http-core/pkg/app"
	"github.com/exgamer/gosdk-http-core/pkg/middleware"
)

func SetRoutes(
	a *app.App,
	cityHandlersFactory *factories.CityHandlersFactory,
) error {
	router, err := http.GetRouter(a)
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
			v1.GET("/city/:id", cityHandlersFactory.CityHandler.View())
		}
	}

	return nil
}
