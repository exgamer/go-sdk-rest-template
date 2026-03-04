package city

import (
	"context"
	"github.com/exgamer/gosdk-core/pkg/debug"
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

type Service struct {
	repository     Repository
	httpRepository HttpRepository
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

	return s.repository.GetById(ctx, id)
}

func (s *Service) Create(ctx context.Context, model *City) (*City, error) {
	model, err := s.repository.Create(ctx, model)

	if err != nil {
		return nil, err
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
