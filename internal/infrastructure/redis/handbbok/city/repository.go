package city

import (
	"context"
	"fmt"
	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"
	rediscore "github.com/exgamer/gosdk-redis-core/pkg/redis"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

type RedisRepository struct {
	client *redis.Client
}

// SetCity - Кеширует город
func (repo *RedisRepository) SetCity(ctx context.Context, model *city.City) error {
	err := rediscore.NewHelper[city.City](ctx, repo.client).SetStruct("service:city:"+strconv.FormatUint(uint64(model.ID), 10), *model, 2*time.Hour)

	if err != nil {
		return err
	}

	return nil
}

// GetCityById - Возвращает город
func (repo *RedisRepository) GetCityById(ctx context.Context, id uint) (*city.City, error) {
	result, err := rediscore.NewHelper[city.City](ctx, repo.client).GetStruct("service:city:" + strconv.FormatUint(uint64(id), 10))

	if err != nil {
		return nil, err
	}

	return result, nil
}

// InvalidateTariffCache Сброс кеша для города
func (repo *RedisRepository) InvalidateTariffCache(ctx context.Context, id uint) error {
	idStr := strconv.FormatUint(uint64(id), 10)

	if err := repo.client.Unlink(
		ctx,
		"service:city:"+idStr,
	).Err(); err != nil {
		return fmt.Errorf("failed to invalidate city cache: %w", err)
	}

	return nil
}
