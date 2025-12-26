package services

import (
	"context"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/database/dto"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/database/models"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/database/repositories"
	"github.com/exgamer/gosdk-core/pkg/debug"
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

func NewCityService(cityRepository *repositories.CityRepository) *CityService {
	return &CityService{
		cityRepository: cityRepository,
	}
}

type CityService struct {
	cityRepository *repositories.CityRepository
}

func (s *CityService) Paginated(ctx context.Context, searchDto *dto.CitySearch) (*pagination.Paginated[models.City], error) {
	paginated, err := s.cityRepository.Paginated(ctx, searchDto)

	if err != nil {
		return nil, err
	}

	return paginated, nil
}

func (s *CityService) GetById(ctx context.Context, id uint) (*models.City, error) {

	if dbg := debug.GetDebugFromContext(ctx); dbg != nil {
		dbg.AddStep("asdfasdf")
	}

	return s.cityRepository.GetById(ctx, id)
}
