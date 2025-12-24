package services

import (
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/models"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/repositories"
)

func NewCityService(cityRepository *repositories.CityRepository) *CityService {
	return &CityService{
		cityRepository: cityRepository,
	}
}

type CityService struct {
	cityRepository *repositories.CityRepository
}

func (s *CityService) GetCity() (*models.City, error) {

	return &models.City{Name: "Алматы"}, nil
}
