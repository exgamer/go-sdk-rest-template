package city

import (
	"context"
	"errors"
	citydomain "github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"
	"github.com/exgamer/gosdk-db-core/pkg/query/helpers"
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
	"gorm.io/gorm"
	"strings"
	"time"
)

func NewPostgresRepository(client *gorm.DB) *PostgresRepository {
	return &PostgresRepository{
		client: client,
	}
}

type PostgresRepository struct {
	client *gorm.DB
}

// Paginated Постраничный список
func (r *PostgresRepository) Paginated(ctx context.Context, searchDto *citydomain.Search) (*pagination.Paginated[citydomain.City], error) {
	helper := helpers.NewGormPaginatedHelper[citydomain.City](ctx, r.client).SetPerPage(searchDto.PerPage)
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
func (r *PostgresRepository) GetById(ctx context.Context, id uint) (*citydomain.City, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var entity citydomain.City
	result := r.client.WithContext(ctx).
		Where("id = ?", id).
		First(&entity)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &entity, nil
}

func (r *PostgresRepository) Create(ctx context.Context, city *citydomain.City) (*citydomain.City, error) {
	model := entityToModel(city)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.client.WithContext(ctx).Create(model)

	if result.Error != nil {
		return nil, result.Error
	}

	return modelToEntity(model), nil
}

func (r *PostgresRepository) Update(ctx context.Context, city *citydomain.City) error {
	model := entityToModel(city)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.client.WithContext(ctx).Save(model)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	result := r.client.WithContext(ctx).Delete(city{}, "id=?", id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresRepository) Activate(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.client.WithContext(ctx).Model(city{}).Where("id=?", id).Update("status", 1)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresRepository) Deactivate(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := r.client.WithContext(ctx).Model(city{}).Where("id=?", id).Update("status", 0)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
