package services

import (
	"context"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/models"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/repositories"
	"github.com/exgamer/gosdk-core/pkg/debug"
)

func NewCityService(cityRepository *repositories.CityRepository) *CityService {
	return &CityService{
		cityRepository: cityRepository,
	}
}

type CityService struct {
	cityRepository *repositories.CityRepository
}

func (s *CityService) GetCity(ctx context.Context, id uint) (*models.City, error) {

	if dbg := debug.GetDebugFromContext(ctx); dbg != nil {
		dbg.AddStep("asdfasdf")
	}

	return s.cityRepository.GetCityById(ctx, id)
}
