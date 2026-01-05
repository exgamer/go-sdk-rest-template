package repositories

import (
	"context"
	"errors"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dal/database/models"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/modules/city/dto"
	"github.com/exgamer/gosdk-db-core/pkg/query/helpers"
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
	"gorm.io/gorm"
	"strings"
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

// Paginated Постраничный список
func (r *CityRepository) Paginated(ctx context.Context, searchDto *dto.CitySearch) (*pagination.Paginated[models.City], error) {
	helper := helpers.NewGormPaginatedHelper[models.City](ctx, r.client).SetPerPage(searchDto.PerPage)
	result, err := helper.Paginated(searchDto.Page, func(client *gorm.DB) *gorm.DB {
		var query []string
		var args []any
		orderBy := "id DESC"

		if searchDto.Id > 0 {
			query = append(query, "city.id=?")
			args = append(args, searchDto.Id)
		}

		return client.
			WithContext(ctx).
			Select("*").
			Where(strings.Join(query, " AND "), args...).
			Order(orderBy)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetById Одна запись по ID
func (r *CityRepository) GetById(ctx context.Context, id uint) (*models.City, error) {
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

func (r *CityRepository) Create(ctx context.Context, model *models.City) (*models.City, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.client.WithContext(ctx).Create(model)

	if result.Error != nil {
		return nil, result.Error
	}

	return model, nil
}

func (r *CityRepository) Update(ctx context.Context, model *models.City) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.client.WithContext(ctx).Save(model)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CityRepository) Delete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	result := r.client.WithContext(ctx).Delete(models.City{}, "id=?", id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CityRepository) Activate(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.client.WithContext(ctx).Model(models.City{}).Where("id=?", id).Update("status", 1)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *CityRepository) Deactivate(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.client.WithContext(ctx).Model(models.City{}).Where("id=?", id).Update("status", 0)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
