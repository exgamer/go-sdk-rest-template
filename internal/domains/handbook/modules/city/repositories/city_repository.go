package repositories

import "github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/models"

func NewCityRepository() *CityRepository {
	return &CityRepository{}
}

type CityRepository struct {
}

func (r *CityRepository) GetCity() (*models.City, error) {

	return &models.City{Name: "Алматы"}, nil
}
