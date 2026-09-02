package city

import (
	"context"
	"github.com/exgamer/gosdk-core/pkg/debug"
	"github.com/exgamer/gosdk-core/pkg/errorreporter"
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

func NewService(repository Repository, cacheRepository CacheRepository) *Service {
	return &Service{
		repository:      repository,
		cacheRepository: cacheRepository,
	}
}

type Service struct {
	repository      Repository
	httpRepository  HttpRepository
	cacheRepository CacheRepository
}

func (s *Service) GetCity(ctx context.Context) (*City, error) {
	return s.httpRepository.GetCity(ctx)
}

func (s *Service) Paginated(ctx context.Context, searchDto *Search) (*pagination.Paginated[City], error) {
	paginated, err := s.repository.Paginated(ctx, searchDto)

	if err != nil {
		return nil, err
	}

	return paginated, nil
}

func (s *Service) GetById(ctx context.Context, id uint) (*City, error) {
	// Так можно добавлять отладочную информацию
	if dbg := debug.GetDebugFromContext(ctx); dbg != nil {
		dbg.AddStep("asdfasdf")
	}

	// Пример errorreporter.CaptureSoft: кеш недоступен (например Redis
	// упал) - запрос не должен из-за этого падать, просто идём в БД.
	// Но факт деградации (кеш не работает) должен долететь до Sentry.
	if s.cacheRepository != nil {
		cached, err := s.cacheRepository.GetCityById(ctx, id)
		if err != nil {
			errorreporter.CaptureSoft(ctx, err, map[string]string{
				"component": "redis_cache",
				"domain":    "city",
			})
		} else if cached != nil {
			return cached, nil
		}
	}

	model, err := s.repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.cacheRepository != nil && model != nil {
		// Не смогли прогреть кеш - тоже не повод валить успешный ответ,
		// просто репортим и едем дальше.
		if err := s.cacheRepository.SetCity(ctx, model); err != nil {
			errorreporter.CaptureSoft(ctx, err, map[string]string{
				"component": "redis_cache",
				"domain":    "city",
			})
		}
	}

	return model, nil
}

func (s *Service) Create(ctx context.Context, model *City) (*City, error) {
	model, err := s.repository.Create(ctx, model)

	if err != nil {
		// Пример errorreporter.CaptureError: ошибку и пробрасываем
		// наверх (как обычно), и одновременно репортим в Sentry - одной
		// строкой. Дедуп внутри errorreporter защитит от повторной
		// отправки, если этот же err ещё раз попадёт в Capture выше по
		// стеку (например в HTTP-транспорте).
		return nil, errorreporter.CaptureError(ctx, err, map[string]string{
			"component": "city_service",
			"operation": "create",
		})
	}

	return model, nil
}

func (s *Service) Update(ctx context.Context, model *City) (*City, error) {
	err := s.repository.Update(ctx, model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *Service) Activate(ctx context.Context, id uint) error {
	err := s.repository.Activate(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Deactivate(ctx context.Context, id uint) error {
	err := s.repository.Deactivate(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	err := s.repository.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
