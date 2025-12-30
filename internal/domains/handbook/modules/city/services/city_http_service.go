package services

import (
	"context"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/http/dto"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/http/repositories"
)

func NewCityHttpService(cityRepository *repositories.CityRepository) *CityHttpService {
	return &CityHttpService{
		cityRepository: cityRepository,
	}
}

type CityHttpService struct {
	cityRepository *repositories.CityRepository
}

func (s *CityHttpService) GetCity(ctx context.Context) (*dto.City, error) {
	return s.cityRepository.GetCity(ctx)
}
