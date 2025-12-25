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
		baseConfig, err := app.GetBaseConfig(a)
		if err != nil {
			return err
		}

		// Middleware
		service.Use(middleware.RequestInfoMiddleware(baseConfig)) // заполнение структур по инфе базовый и http
		service.Use(middleware.LoggerMiddleware())                // форматированные логи

		v1 := service.Group("/v1")
		{
			v1.Use(middleware.FormattedResponseMiddleware()) // форматированный ответ
			v1.Use(middleware.MetricsMiddleware(a))          // метрики прометея
			v1.GET("/city", cityHandlersFactory.CityHandler.GetCity())
		}
	}

	return nil
}
