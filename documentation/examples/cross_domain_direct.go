//go:build ignore

// Случай 1: прямая зависимость между сервисами через интерфейс.
// Пример: OrderService использует CityService для валидации города.

// ===== internal/domains/order/order/service.go =====

package order

import (
	"context"
	"errors"

	"github.com/exgamer/gosdk-core/pkg/exception"
	citydomain "github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"
)

// CityServiceInterface — контракт с сервисом города.
// Реализуется конкретным city.Service из домена handbook/city.
type CityServiceInterface interface {
	GetById(ctx context.Context, id uint) (*citydomain.City, error)
}

type Service struct {
	repository           Repository
	cityServiceInterface CityServiceInterface
}

func NewService(repository Repository, cityServiceInterface CityServiceInterface) *Service {
	return &Service{
		repository:           repository,
		cityServiceInterface: cityServiceInterface,
	}
}

func (s *Service) Create(ctx context.Context, m *Order) (*Order, error) {
	city, err := s.cityServiceInterface.GetById(ctx, m.CityID)
	if err != nil {
		return nil, err
	}

	if city == nil {
		return nil, exception.NewNotFoundException(errors.New("city not found"), false)
	}

	return s.repository.Create(ctx, m)
}

// ===== internal/app/bootstrap/city/module.go =====
// Шаг 1: city.Module регистрирует свой сервис в DI.
// order.Module запускается после и получает его.

// func (m *Module) Init(a *app.App) error {
//     client, _ := postgresDi.GetDefaultPostgresConnection(a.Container)
//
//     repoFactory := newRepositoriesFactory(client)
//     svcFactory  := newServicesFactory(repoFactory)
//
//     // зарегистрировать в DI, чтобы order.Module мог получить
//     di.Register(a.Container, svcFactory.CityService)
//
//     hdlFactory := newHandlersFactory(svcFactory)
//     entrypoint.SetRoutes(a, hdlFactory.CityHandler)
//
//     return nil
// }

// ===== internal/app/bootstrap/order/services_factory.go =====

// func newServicesFactory(repos *repositoriesFactory, cityServiceInterface *city.Service) *servicesFactory {
//     return &servicesFactory{
//         OrderService: order.NewService(repos.PostgresRepository, cityServiceInterface),
//     }
// }

// ===== internal/app/bootstrap/order/module.go =====

// func (m *Module) Init(a *app.App) error {
//     client, _       := postgresDi.GetDefaultPostgresConnection(a.Container)
//     cityService, _  := di.Resolve[*city.Service](a.Container)
//
//     repoFactory := newRepositoriesFactory(client)
//     svcFactory  := newServicesFactory(repoFactory, cityService)
//     // ...
//     return nil
// }

// ===== internal/app/app.go — порядок регистрации модулей =====

// err = appInstance.RegisterAndInitModules(
//     &city.Module{},   // ← сначала: регистрирует CityService в DI
//     &order.Module{},  // ← потом: получает CityService из DI
// )
