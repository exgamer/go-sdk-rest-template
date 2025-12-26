package repositories

import (
	"context"
	"errors"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/models"
	"gorm.io/gorm"
	"time"
)

func NewCityRepository(client *gorm.DB) *CityRepository {
	return &CityRepository{
		client: client,
	}
}

type CityRepository struct {
	client *gorm.DB
}

func (r *CityRepository) GetCityById(ctx context.Context, id uint) (*models.City, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var model models.City
	result := r.client.WithContext(ctx).
		Where("id = ?", id).
		First(&model)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &model, nil
}
